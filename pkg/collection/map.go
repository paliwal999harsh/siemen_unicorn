package collection

type Map[K any, V any] interface {
	Get(K) (V, bool)
	Put(K, V)
	Remove(K)
	Size() int
	Empty() bool
}

type nativeMap[K comparable, V any] struct {
	data map[K]V
}

func (n *nativeMap[K, V]) Get(key K) (V, bool) {
	var zero V
	if n.Empty() {
		return zero, false
	}
	if v, ok := n.data[key]; ok {
		return v, ok
	}
	return zero, false
}

func (n *nativeMap[K, V]) Put(key K, value V) {
	n.data[key] = value
}

func (n *nativeMap[K, V]) Remove(key K) {
	delete(n.data, key)
}

func (n *nativeMap[K, V]) Size() int {
	return len(n.data)
}

func (n *nativeMap[K, V]) Empty() bool {
	return len(n.data) == 0
}

func NewNativeMap[K comparable, V any]() Map[K, V] {
	return &nativeMap[K, V]{
		data: make(map[K]V),
	}
}
