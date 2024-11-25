package deque

const minCapacity = 16

type Deque[T any] struct {
	buf                       []T
	head, tail, count, minCap int
}

func (q *Deque[T]) Cap() int {
	if q == nil {
		return 0
	}
	return len(q.buf)
}

func (q *Deque[T]) Len() int {
	if q == nil {
		return 0
	}
	return q.count
}

func (q *Deque[T]) PushBack(item T) {
	q.growIfFull()

	q.buf[q.tail] = item
	q.tail = q.next(q.tail)
	q.count++
}

func (q *Deque[T]) PushFront(item T) {
	q.growIfFull()

	q.head = q.prev(q.head)
	q.buf[q.head] = item
	q.count++
}

func (q *Deque[T]) PopFront() T {
	if q.count <= 0 {
		panic("empty queue")
	}

	var defaultTVal T

	ret := q.buf[q.head]
	q.buf[q.head] = defaultTVal
	q.head = q.next(q.head)
	q.count--

	q.shrink()
	return ret
}

func (q *Deque[T]) PopBack() T {
	if q.count <= 0 {
		panic("empty queue")
	}

	var defaultTVal T
	ret := q.buf[q.tail]
	q.buf[q.tail] = defaultTVal
	q.count--

	q.shrink()
	return ret
}

func (q *Deque[T]) Clear() {
	var defaultTVal T
	modBits := len(q.buf) - 1
	h := q.head
	for i := 0; i < q.Len(); i++ {
		q.buf[(h+i)&modBits] = defaultTVal
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

func (q *Deque[T]) Grow(size int) {
	if size < 0 {
		panic("negative grow count")
	}
	c := q.Cap()
	l := q.Len()

	if size <= c-l {
		return
	}

	if c == 0 {
		c = minCapacity
	}

	newSize := l + size
	for c < newSize {
		c <<= 1
	}

	if l == 0 {
		q.buf = make([]T, c)
		q.head = 0
		q.tail = 0
	} else {
		q.resize(c)
	}
}

func (q *Deque[T]) SetBaseCap(baseCap int) {
	minCap := minCapacity
	for minCap < baseCap {
		minCap <<= 1
	}
	q.minCap = minCap
}

func (q *Deque[T]) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1)
}

func (q *Deque[T]) next(i int) int {
	return (i + 1) & (len(q.buf) - 1)
}

func (q *Deque[T]) growIfFull() {
	if q.count != len(q.buf) {
		return
	}

	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = q.minCap
		}
		q.buf = make([]T, q.minCap)
		return
	}
	q.resize(q.count << 1)
}

func (q *Deque[T]) shrink() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize(q.count << 1)
	}
}

func (q *Deque[T]) resize(newSize int) {
	newBuf := make([]T, newSize)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}
