package apiv1

import "container/list"

type Queue struct {
	list *list.List
}

func (q *Queue) InitQueue() {
	q.list = list.New()
}

func (q *Queue) Enqueue(v int) {
	q.list.PushBack(v)
}

func (q *Queue) Dequeue() {
	q.list.Remove(q.list.Front())
}

func (q *Queue) Head() int {
	return q.list.Front().Value.(int)
}

func (q *Queue) Size() int {
	return q.list.Len()
}
