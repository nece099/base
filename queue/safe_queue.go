package queue

import "sync"

type Queue struct {
	ch     chan interface{}
	closed bool
	mutex  *sync.Mutex
}

func NewQueue(len int) *Queue {
	return &Queue{
		ch:     make(chan interface{}, len),
		closed: false,
		mutex:  &sync.Mutex{},
	}
}

func (q *Queue) Push(d interface{}) {
	q.ch <- d
}

func (q *Queue) Pop() (interface{}, bool) {
	d, ok := <-q.ch
	return d, ok
}

func (q *Queue) Len() int {
	return len(q.ch)
}

func (q *Queue) Close() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if !q.closed {
		close(q.ch)
		q.closed = true
	}
}

func (q *Queue) IsClosed() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.closed
}
