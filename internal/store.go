package internal

type GenericMap[K comparable, V any] struct {
	M map[K]V
}

func (m* GenericMap[K,V]) Get(key K) (V, bool) {
	value, ok := m.M[key]
	return value, ok
}

func (m *GenericMap[K, V]) Put(key K, value V) {
	m.M[key] = value;
}

func (m *GenericMap[K, V]) Delete(key K) {
		delete(m.M, key)
}

func (m *GenericMap[K, V]) Len() int {
	return len(m.M)
}

func (m *GenericMap[K, V]) All(f func(key K, val V) bool) {
	for k,v := range m.M {
		if !f(k,v) {
			break;
		}
	}
}	


type HashTable[K comparable, V any] interface {
	Put()
	Get(key K) (V, bool)
	Delete(key K)
	Len() int
	All(f func(key K, val V) bool)
}