package chronotable

import "sync"

type ChronoTable struct {
	M   map[string]interface{}
	mtx sync.RWMutex
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

}

func (m *ChronoTable) Delete(key string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.M, key)
}

func (m *ChronoTable) Len() int {
	return len(m.M)
}


func New() *ChronoTable {
	return &ChronoTable{
		M: make(map[string]interface{}),
	}
}
