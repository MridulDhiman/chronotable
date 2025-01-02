package chronotable

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"github.com/MridulDhiman/chronotable/config"
	"github.com/MridulDhiman/chronotable/internal/aof"
	"github.com/MridulDhiman/chronotable/internal/decoder"
	"github.com/MridulDhiman/chronotable/internal/snapshot"
)

type ChronoTable struct {
	// TODO: make it private
	M        map[string]interface{}
	mtx      sync.RWMutex
	aof      *aof.AOF
	snapshot *snapshot.SnapShot
	// TODO: make it private
	ConfigHandler *config.ConfigHandler
	logger        *log.Logger
}


func New(opts *Options) *ChronoTable {
	t := &ChronoTable{
		M:             make(map[string]interface{}),
		ConfigHandler: config.NewConfigHandler(opts.Mode),
		logger:        opts.Logger,
	}

	if !opts.Initialized {
		go t.ConfigHandler.UpdateConfigFile(0)
	}

	if opts.EnableAOF {
		t.aof = aof.New(opts.AOFPath, opts.Initialized)
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

func (m *ChronoTable) Put(key string, value interface{}, options ...InputOpts) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	inputOpts:= m.handleInputOpts(options...)
	if m.aof != nil && !inputOpts.IsReplayed {
		m.aof.Log(aof.MustFormat(key, value, aof.PutOp))
	}
	m.M[key] = value
}

func (m *ChronoTable) Delete(key string, options ...InputOpts) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	inputOpts:= m.handleInputOpts(options...)
	if m.aof != nil && !inputOpts.IsReplayed {
		m.aof.Log(aof.MustFormat(key, nil, aof.DeleteOp))
	}
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

// if snapshot configured in chronotable,
// creates new snapshot from current state of the hash table, 
// clears the log, else returns nil 
func (m *ChronoTable) Commit() *snapshot.Version {
	if m.snapshotEnabled() {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		newSnapShot, err := m.snapshot.Create(m.M, m.ConfigHandler)
		if err != nil {
			fmt.Println("Error in creating snapshot: ", err)
			return nil
		}
		// TODO: handling properly if snapshot creation successful, but could not clear the log 
		if m.aofEnabled() {
			if err := m.aof.Clear(); err != nil {
					fmt.Println("Error in clearing log: ", err)
					return nil
			}
		}
		return newSnapShot
	}
	return nil
}

func (m *ChronoTable) Timetravel(version int64) {
	if m.snapshotEnabled() {
		desiredVersion, ok := m.snapshot.GetVersion(version)
		if ok {
			m.Clear()
			m.Copy(desiredVersion.Data)
			m.snapshot.CurrentVersion = version
			m.ConfigHandler.UpdateConfigFile(version)
		}
	}
}

func (m *ChronoTable) List() {
	for k, v := range m.M {
		fmt.Printf("Key: %s, Value: %v\n", k, v)
	}
}


// Fetch Latest Snapshot file and return the desired version
func (m *ChronoTable) ReplayOnRestart(currentVersion, latestVersion int64) error {
	if latestVersion != int64(0) {
		snapshotFile := strconv.FormatInt(currentVersion, 10) + config.SNAPSHOT_EXT
		file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, snapshotFile), os.O_RDONLY, os.FileMode(0644))
		if err != nil {
			return fmt.Errorf("(error) could not open file: %v", err)
		}
		defer file.Close()
		// TODO: concretions object creation outside function
		binDecoder := decoder.NewDecoder(file)
		desiredVersion := new(snapshot.Version)
		if err := binDecoder.Decode(desiredVersion); err != nil {
			return err
		}
		m.Clear()
		m.snapshot.SetCurrentVersion(desiredVersion.Id)
		m.snapshot.SetLatestVersion(latestVersion)
		m.Copy(desiredVersion.Data)
	}

	m.populate(m.aof.MustReplay())
	return nil
}

// ~~~ PRIVATE METHODS ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
func (m *ChronoTable) snapshotEnabled() bool {
	return m.snapshot != nil;
}

func (m *ChronoTable) aofEnabled() bool {
	return m.aof != nil;
}

func (m * ChronoTable) handleInputOpts(inputOpts ...InputOpts) InputOpts {
	var isReplayed bool = false;
	if len(inputOpts) == 1  {
		isReplayed = inputOpts[0].IsReplayed
	}
return InputOpts{
	IsReplayed: isReplayed,
}
}
func (m * ChronoTable) populate(operations []*aof.Operation) {

	for _, operation := range operations {
		if operation.OperationType == aof.PutOp {
			m.Put(operation.Key, operation.Value, InputOpts{
				IsReplayed: true,
			})
		} else if operation.OperationType == aof.DeleteOp {
			m.Delete(operation.Key, InputOpts{
				IsReplayed: true,
			})
		}
	}
}
