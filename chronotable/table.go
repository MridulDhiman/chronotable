package chronotable

import (
	"fmt"
	"sync"

	"github.com/MridulDhiman/chronotable/internal/aof"
	"github.com/MridulDhiman/chronotable/internal/snapshot"
)



type ChronoTable struct {
	M        map[string]interface{}
	mtx      sync.RWMutex
	aof      *aof.AOF
	snapshot *snapshot.SnapShot
}

func New(opts *Options) *ChronoTable {
	t := &ChronoTable{
		M: make(map[string]interface{}),
	}

	if opts.EnableAOF {
		t.aof = aof.New(opts.AOFPath)
	}

	if opts.EnableSnapshot {
		t.snapshot = snapshot.New()
	}

	return t
}

func (m *ChronoTable) Get(key string) (interface{}, bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	value, ok := m.M[key]
	return value, ok
}

func (m *ChronoTable) Put(key string, value interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.M[key] = value
	if m.aof != nil {
	 m.aof.Log(aof.Format(key, value))
	}
}

func (m *ChronoTable) Delete(key string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.M, key)
}

func (m *ChronoTable) Len() int {
	return len(m.M)
}

func (m *ChronoTable) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for k := range m.M {
		delete(m.M, k)
	}
}

func (m *ChronoTable) Copy(m2 map[string]any) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for k, v := range m2 {
		m.M[k] = v
	}
}

func (m *ChronoTable) Commit() *snapshot.Version {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	latestVersion, ok := m.snapshot.GetLatestVersion()
	if !ok {
		return nil
	}
	var AOFStart int64 = 0
	if latestVersion != nil {
		AOFStart = latestVersion.AOFEnd
	}

	newSnapShot, err := m.snapshot.Create(m.M, AOFStart, m.aof.SeekCurrent)
	if err != nil {
		fmt.Println("Error in creating snapshot: ", err)
		return nil
	}

	return newSnapShot
}

func (m *ChronoTable) SnapshotEnabled() bool {
	return m.snapshot == nil
}

func (m *ChronoTable) Timetravel(version int) {
	desiredVersion, ok := m.snapshot.GetVersion(version)
	if ok {
		m.Clear()
		m.Copy(desiredVersion.Data)
		m.snapshot.CurrentVersion = version
	}
}

// get the current version's new insertions 
func (m *ChronoTable) ChangesCurrent() {
	 desiredVersion, _ := m.snapshot.GetVersion(m.snapshot.CurrentVersion); 
	m.aof.Replay(desiredVersion.AOFStart, desiredVersion.AOFEnd)
}

// get the changes till current version
func (m *ChronoTable) ChangesTill() {
	desiredVersion, _ := m.snapshot.GetVersion(m.snapshot.CurrentVersion); 
   m.aof.Replay(0, desiredVersion.AOFEnd)
}
