package ordcol

import (
	"container/heap"
	"container/list"

	"golang.org/x/exp/constraints"
)

type MinHeap[K constraints.Ordered, V any] []*list.Element

func (h MinHeap[K, V]) Len() int {
	return len(h)
}
func (h MinHeap[K, V]) Less(i, j int) bool {
	return h[i].Value.(Element[K, V]).key < h[j].Value.(Element[K, V]).key
}
func (h MinHeap[K, V]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap[K, V]) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*list.Element))
}

func (h *MinHeap[K, V]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type CollectionImpl[K constraints.Ordered, V any] struct {
	l *list.List
	m map[K]*list.Element
	h *MinHeap[K, V]
}

type Element[K constraints.Ordered, V any] struct {
	key   K
	value V
}

type IteratorImpl[K constraints.Ordered, V any] struct {
	l     *list.List
	elem  *list.Element
	order IterationOrder
}

func (i *IteratorImpl[K, V]) HasNext() bool {
	return i.elem != nil
}

// Next returns ErrEmptyIterator if there is no more elements (HasNext() is false)
func (i *IteratorImpl[K, V]) Next() (K, V, error) {
	if !i.HasNext() {
		return *new(K), *new(V), ErrEmptyIterator
	}
	elem := i.elem.Value.(Element[K, V])
	if i.order == ByInsertion {
		i.elem = i.elem.Next()
	} else {
		i.elem = i.elem.Prev()
	}
	return elem.key, elem.value, nil
}

// Add adds specified value to the collection and associates it with the specified key.
//
// If there already is a value associated with the same key, ErrDuplicateKey is returned.
func (c *CollectionImpl[K, V]) Add(key K, value V) error {
	_, found := c.m[key]
	if found {
		return ErrDuplicateKey
	}
	elem := c.l.PushBack(Element[K, V]{key, value})
	heap.Push(c.h, elem)
	c.m[key] = elem
	return nil
}

// DelMin removes the value associated with the minimum key from the collection and returns it.
//
// If the collection is empty, ErrEmptyCollection is returned.
func (c *CollectionImpl[K, V]) DelMin() (K, V, error) {
	if c.Len() == 0 {
		return *new(K), *new(V), ErrEmptyCollection
	}

	elem := heap.Pop(c.h).(*list.Element)
	kv := c.l.Remove(elem).(Element[K, V])
	delete(c.m, kv.key)
	return kv.key, kv.value, nil
}

// Len returns the number of elements in the collection.
func (c *CollectionImpl[K, V]) Len() int {
	return len(c.m)
}

// IterateBy returns an iterator over the elements of the collection in the specified order.
// Any modification of the collection corrupts all the iterators obtained prior the modification.
//
// If invalid order is passed, the function panics with ErrUnknownOrder
func (c *CollectionImpl[K, V]) IterateBy(order IterationOrder) Iterator[K, V] {
	if order == ByInsertion {
		return &IteratorImpl[K, V]{
			c.l, c.l.Front(), ByInsertion,
		}
	} else if order == ByInsertionRev {
		return &IteratorImpl[K, V]{
			c.l, c.l.Back(), ByInsertionRev,
		}
	} else {
		panic(ErrUnknownOrder)
	}
}

// At returns the value associated with the specified key.
func (c *CollectionImpl[K, V]) At(key K) (V, bool) {
	val, found := c.m[key]
	if !found {
		return *new(V), false
	}
	return val.Value.(Element[K, V]).value, true
}

func NewCollection[K constraints.Ordered, V any]() Collection[K, V] {
	h := make(MinHeap[K, V], 0)
	heap.Init(&h)

	return &CollectionImpl[K, V]{
		l: list.New(),
		m: make(map[K]*list.Element),
		h: &h,
	}
}
