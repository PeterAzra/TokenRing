package deque

import "sync"

type CDeque[T any] struct {
	innerQ Deque[T]
	mu     sync.Mutex
}

func (q *CDeque[T]) Cap() int {
	return q.innerQ.Cap()
}

func (q *CDeque[T]) Len() int {
	return q.innerQ.Len()
}

func (q *CDeque[T]) PushBack(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.innerQ.PushBack(item)
}

func (q *CDeque[T]) PushFront(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.innerQ.PushFront(item)
}

func (q *CDeque[T]) PopFront() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.innerQ.PopFront()
}

func (q *CDeque[T]) PopBack() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.innerQ.PopBack()
}

func (q *CDeque[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.innerQ.Clear()
}

func (q *CDeque[T]) Grow(size int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.innerQ.Grow(size)
}

func (q *CDeque[T]) SetBaseCap(baseCap int) {
	q.innerQ.SetBaseCap(baseCap)
}
