package db

import "fmt"

type ErrDataNotFound[T any] struct {
	id T
}

func NewErrDataNotFound[T any](id T) *ErrDataNotFound[T] {
	return &ErrDataNotFound[T]{id}
}

func (e *ErrDataNotFound[T]) Error() string {
	return fmt.Sprintf("Data cound not be found by identifier (%v)", e.id)
}
