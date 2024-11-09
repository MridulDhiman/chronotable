package chronotable

import (
	"sync"

	"github.com/MridulDhiman/chronotable/internal/aof"
)

type ChronoTable struct {
	M   map[string]interface{}
	mtx sync.RWMutex
	aof *aof.AOF
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


func New(opts *Options) *ChronoTable {
	t:= &ChronoTable{
		M: make(map[string]interface{}),
	}

	if opts.EnableAOF {
		t.aof = aof.New(opts.AOFPath)
	}

	return t;
}
