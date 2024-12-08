package chronotable

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"github.com/MridulDhiman/chronotable/config"
	"github.com/MridulDhiman/chronotable/internal/aof"
	"github.com/MridulDhiman/chronotable/internal/decoder"
	"github.com/MridulDhiman/chronotable/internal/snapshot"
	"github.com/MridulDhiman/chronotable/internal/utils"
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
		t.aof = aof.New(opts.AOFPath, opts.Initialized)
	}

	if opts.EnableSnapshot {
		t.snapshot = snapshot.New(opts.Initialized)
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
		fmt.Println("(error) could not get latest version")
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
		utils.UpdateConfigFile(version)
	}
}

// get the current version's new insertions
func (m *ChronoTable) ChangesCurrent() {
	desiredVersion, _ := m.snapshot.GetVersion(m.snapshot.CurrentVersion)
	m.aof.Replay(desiredVersion.AOFStart, desiredVersion.AOFEnd)
}

func (m *ChronoTable) List() {
	for k,v := range m.M {
	fmt.Printf("Key: %s, Value: %v\n", k, v)
	}
}

// get the changes till current version
func (m *ChronoTable) ChangesTill() {
	desiredVersion, _ := m.snapshot.GetVersion(m.snapshot.CurrentVersion)
	m.aof.Replay(0, desiredVersion.AOFEnd)
}

// Fetch Latest Snapshot file and return the desired version
func (m *ChronoTable) ReplayOnRestart(currentVersion int, latestVersion int) error {
	snapshotFile := strconv.Itoa(currentVersion) + config.SNAPSHOT_EXT
	file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, snapshotFile), os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("(error) could not open file: %v", err)
	}
	defer file.Close()
	binDecoder := decoder.NewDecoder(file)
	desiredVersion := new(snapshot.Version)
	if err := binDecoder.Decode(desiredVersion); err != nil {
		return err
	}
	m.Clear()
	m.snapshot.SetCurrentVersion(desiredVersion.Id)
	m.snapshot.SetLatestVersion(latestVersion)
	m.Copy(desiredVersion.Data)
	return nil
}
