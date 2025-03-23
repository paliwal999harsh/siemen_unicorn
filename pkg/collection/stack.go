package collection

import "errors"

type Stack[T any] interface {
	Push(T) (bool, error)
	Pop() (T, error)
	Peek() (T, error)
	Size() int
	Empty() bool
}

type sliceStack[T any] struct {
	list []T
}

func (s *sliceStack[T]) Push(item T) (bool, error) {
	s.list = append(s.list, item)
	return true, nil
}

func (s *sliceStack[T]) Pop() (T, error) {
	if s.Empty() {
		var zero T
		return zero, errors.New("empty stack")
	}
	item := s.list[s.Size()-1]
	s.list = s.list[:s.Size()-1]
	return item, nil
}

func (s *sliceStack[T]) Peek() (T, error) {
	return s.list[0], nil
}

func (s *sliceStack[T]) Size() int {
	return len(s.list)
}

func (s *sliceStack[T]) Empty() bool {
	return len(s.list) == 0
}

func NewSliceStack[T any]() Stack[T] {
	return &sliceStack[T]{}
}
