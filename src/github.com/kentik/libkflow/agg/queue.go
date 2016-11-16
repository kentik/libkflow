package agg

import (
	"fmt"

	"github.com/kentik/libkflow/chf"
)

/**
Stolen and hacked up from https://github.com/eapache/queue/blob/master/queue.go
*/

const minQueueLen = 1028

// Queue represents a single instance of the queue data structure.
type Queue struct {
	buf                             []*chf.CHF
	head, tail, count, max, dropped int
}

// New constructs and returns a new Queue.
func New(max int) *Queue {
	return &Queue{
		buf: make([]*chf.CHF, minQueueLen),
		max: max,
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newBuf := make([]*chf.CHF, q.count*2)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		copy(newBuf, q.buf[q.head:len(q.buf)])
		copy(newBuf[len(q.buf)-q.head:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *Queue) Enqueue(elem *chf.CHF) error {

	if q.count >= q.max {
		q.dropped++
		return fmt.Errorf("Max q depth exceeeded: %d", q.max)
	}

	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
	return nil
}

// Returns oldest i elements in q. if len < i, returns len
func (q *Queue) Dequeue(i int, overflow int) ([]*chf.CHF, int, float32) {

	get := i

	// If we are full, return overflow situation
	if q.count >= q.max {
		get = overflow
	}

	// But bound by amount actually in q.
	if i >= q.count {
		get = q.count
	}

	// Bad inputs get nil output.
	if get <= 0 {
		return nil, 0, 0
	}

	qLen := len(q.buf)
	res := make([]*chf.CHF, get)
	nilSlice := make([]*chf.CHF, get)

	if q.head+get < qLen {
		copy(res, q.buf[q.head:q.head+get])
		// Nil out buffer too.
		copy(q.buf[q.head:q.head+get], nilSlice)
	} else {
		// Have to do it in two phases
		tPt := qLen - q.head
		copy(res, q.buf[q.head:qLen])
		copy(res[tPt:], q.buf[0:get-tPt])

		// Nil out buffer too, down here
		copy(q.buf[q.head:qLen], nilSlice[0:tPt])
		copy(q.buf[0:get-tPt], nilSlice[0:get-tPt])
	}

	q.head = (q.head + get) % len(q.buf)
	q.count = q.count - get
	if len(q.buf) > minQueueLen && q.count*4 == len(q.buf) {
		q.resize()
	}

	// Calculate a rate adjust here.
	rateAdj := (float32(get + q.dropped)) / float32(get)
	q.dropped = 0

	return res, get, rateAdj
}
