package main

import (
	"errors"
)

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:] // Slice off the first element
	return item, nil
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}
