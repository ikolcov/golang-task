package ordcol

import "golang.org/x/exp/constraints"

type CollectionImpl[K constraints.Ordered, V any] struct{}

// Add adds specified value to the collection and associates it with the specified key.
//
// If there already is a value associated with the same key, ErrDuplicateKey is returned.
func (c *CollectionImpl[K, V]) Add(key K, value V) error {
	panic("not implemented")
}

// DelMin removes the value associated with the minimum key from the collection and returns it.
//
// If the collection is empty, ErrEmptyCollection is returned.
func (c *CollectionImpl[K, V]) DelMin() (K, V, error) {
	panic("not implemented")
}

// Len returns the number of elements in the collection.
func (c *CollectionImpl[K, V]) Len() int {
	panic("not implemented")
}

// IterateBy returns an iterator over the elements of the collection in the specified order.
// Any modification of the collection corrupts all the iterators obtained prior the modification.
//
// If invalid order is passed, the function panics with ErrUnknownOrder
func (c *CollectionImpl[K, V]) IterateBy(order IterationOrder) Iterator[K, V] {
	panic("not implemented")
}

// At returns the value associated with the specified key.
func (c *CollectionImpl[K, V]) At(key K) (V, bool) {
	panic("not implemented")
}

func NewCollection[K constraints.Ordered, V any]() Collection[K, V] {
	return &CollectionImpl[K, V]{}
}
