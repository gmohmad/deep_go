package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type cmpFn[K comparable] func(a, b K) int

type Node[K comparable, V any] struct {
	left, right *Node[K, V]
	key         K
	value       V
}

type OrderedMap[K comparable, V any] struct {
	root *Node[K, V]
	size int
	cmp  cmpFn[K]
}

func NewOrderedMap[K comparable, V any](cmp cmpFn[K]) OrderedMap[K, V] {
	return OrderedMap[K, V]{cmp: cmp}
}

func NewNode[K comparable, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{key: key, value: value}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	curr := m.root

	var prev *Node[K, V]
	for curr != nil {
		prev = curr
		switch {
		case m.cmp(prev.key, key) < 0:
			curr = curr.right
		case m.cmp(prev.key, key) > 0:
			curr = curr.left
		default:
			curr.value = value
			return
		}
	}

	newNode := NewNode(key, value)
	switch {
	case prev == nil:
		m.root = newNode
	case m.cmp(prev.key, key) < 0:
		prev.right = newNode
	case m.cmp(prev.key, key) > 0:
		prev.left = newNode
	}
	m.size++
}

func (m *OrderedMap[K, V]) Erase(key K) {
	if newRoot, deleted := deleteNode(m.root, key, m.cmp); deleted {
		m.size--
		m.root = newRoot
	}
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	curr := m.root

	for curr != nil {
		switch {
		case m.cmp(curr.key, key) < 0:
			curr = curr.right
		case m.cmp(curr.key, key) > 0:
			curr = curr.left
		default:
			return true
		}
	}
	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	traverse(m.root, action)
}

func deleteNode[K comparable, V any](node *Node[K, V], key K, cmp cmpFn[K]) (*Node[K, V], bool) {
	if node == nil {
		return nil, false
	}

	var deleted bool
	switch {
	case cmp(node.key, key) < 0:
		node.right, deleted = deleteNode(node.right, key, cmp)
	case cmp(node.key, key) > 0:
		node.left, deleted = deleteNode(node.left, key, cmp)
	default:
		deleted = true
		if node.left == nil {
			return node.right, deleted
		}
		if node.right == nil {
			return node.left, deleted
		}

		rightMin := findMin(node.right)
		node.key = rightMin.key
		node.value = rightMin.value
		node.right, _ = deleteNode(node.right, rightMin.key, cmp)
	}

	return node, deleted
}

func traverse[K comparable, V any](node *Node[K, V], action func(K, V)) {
	if node == nil {
		return
	}
	traverse(node.left, action)
	action(node.key, node.value)
	traverse(node.right, action)
}

func findMin[K comparable, V any](n *Node[K, V]) *Node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int](cmp.Compare)
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
