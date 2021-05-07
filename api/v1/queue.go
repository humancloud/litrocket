package apiv1

import "sync"

var rwmutex sync.RWMutex

type Queue struct {
	data []byte
	next *Queue
}

func NewQueue() *Queue {
	queue := new(Queue)
	queue.next = nil
	return queue
}

func (q *Queue) Dequeue() {
	rwmutex.Lock()
	var node = q.next
	if node != nil {
		q.next = node.next
	}
	rwmutex.Unlock()
}

func (q *Queue) Enqueue(data []byte) {
	var walk *Queue

	node := new(Queue)
	node.data = make([]byte, len(data))
	copy(node.data, data)
	node.next = nil

	rwmutex.RLock()
	for walk = q; walk.next != nil; walk = walk.next {
	}
	rwmutex.RUnlock()

	rwmutex.Lock()
	walk.next = node
	rwmutex.Unlock()
}

func (q *Queue) Head() []byte {
	rwmutex.RLock()
	if q.next != nil {
		rwmutex.RUnlock()
		return q.next.data
	}
	rwmutex.RUnlock()
	return []byte("")
}

func (q *Queue) IsEmpty() bool {
	return q.next == nil
}
