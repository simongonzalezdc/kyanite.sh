package ui

// GenericList is a type-safe generic list implementation
type GenericList[T any] struct {
	items []T
}

// NewGenericList creates a new generic list
func NewGenericList[T any]() *GenericList[T] {
	return &GenericList[T]{
		items: make([]T, 0),
	}
}

// Add adds an item to the list
func (l *GenericList[T]) Add(item T) {
	l.items = append(l.items, item)
}

// Get gets an item at the given index
func (l *GenericList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, false
	}
	return l.items[index], true
}

// Remove removes an item at the given index
func (l *GenericList[T]) Remove(index int) bool {
	if index < 0 || index >= len(l.items) {
		return false
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
	return true
}

// Len returns the number of items
func (l *GenericList[T]) Len() int {
	return len(l.items)
}

// Clear removes all items
func (l *GenericList[T]) Clear() {
	l.items = make([]T, 0)
}

// Items returns all items as a slice
func (l *GenericList[T]) Items() []T {
	return l.items
}

// Filter returns a new list with items that match the predicate
func (l *GenericList[T]) Filter(predicate func(T) bool) *GenericList[T] {
	result := NewGenericList[T]()
	for _, item := range l.items {
		if predicate(item) {
			result.Add(item)
		}
	}
	return result
}

// Map transforms the list using the given function
func (l *GenericList[T]) Map(transform func(T) T) *GenericList[T] {
	result := NewGenericList[T]()
	for _, item := range l.items {
		result.Add(transform(item))
	}
	return result
}
