package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Integer interface {
	int8 | int16 | int32 | int64 | int
}

type CircularQueue[T Integer] struct {
	values []T
	back   int
	front  int
}

func NewCircularQueue[T Integer](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
		back:   -1,
		front:  -1,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.back = (q.back + 1) % len(q.values)
	q.values[q.back] = value

	if q.front == -1 {
		q.front = 0
	}

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.values[q.front] = 0

	if q.front == 0 && q.back == 0 {
		q.front--
		q.back--
	} else {
		q.front = (q.front + 1) % len(q.values)
	}

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.front]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.back]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.front == -1 && q.back == -1
}

func (q *CircularQueue[T]) Full() bool {
	return (q.front == 0 && q.back == len(q.values)-1) || (q.front == q.back+1)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
