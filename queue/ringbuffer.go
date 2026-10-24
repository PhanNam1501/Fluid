package queue

const minQueuLen = 16

type Queue[V any] struct {
	buf                []*V
	head, tail, length int
}

func NewQueue[V any]() *Queue[V] {
	return &Queue[V]{
		buf: make([]*V, minQueuLen),
	}
}

func (q *Queue[V]) resize() {
	newBuf := make([]*V, q.length<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		copy(newBuf, q.buf[q.head:])
		copy(newBuf, q.buf[:q.tail])
	}
	q.buf = newBuf
	q.head = 0
	q.tail = q.length
}

func (q *Queue[V]) Length() int {
	return q.length
}

func (q *Queue[V]) Head() *V {
	if q.length == 0 {
		panic("queue: Head() called on empty queue")
	}
	return q.buf[q.head]
}

func (q *Queue[V]) Push(elem V) {
	if q.length == len(q.buf) {
		q.resize()
	}
	q.buf[q.tail] = &elem
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.length++
}

func (q *Queue[V]) Pop() {
	if q.length == 0 {
		panic("queue: Remove() called on empty queue")
	}
	q.buf[q.head] = nil
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.length--

	if len(q.buf) > minQueuLen && (q.length<<2) == len(q.buf) {
		q.resize()
	}
}
