package collection

import (
	"errors"
)

type Queue[T any] interface {
	Offer(T) (bool, error)
	Poll() (T, error)
	Size() int
	Empty() bool
}

type sliceQueue[T any] struct {
	list []T
}

func (q *sliceQueue[T]) Offer(item T) (bool, error) {
	q.list = append(q.list, item)
	return true, nil
}

func (q *sliceQueue[T]) Poll() (T, error) {
	if len(q.list) == 0 {
		var zero T
		return zero, errors.New("empty queue")
	}
	item := q.list[0]
	q.list = q.list[1:]
	return item, nil
}

func (q *sliceQueue[T]) Size() int {
	return len(q.list)
}

func (q *sliceQueue[T]) Empty() bool {
	return len(q.list) == 0
}

func NewSliceQueue[T any]() Queue[T] {
	return &sliceQueue[T]{}
}
