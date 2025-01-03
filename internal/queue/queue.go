package queue

import "sync"

type Queue struct {
	messages []string
	mutex    sync.Mutex
}

// NewQueue create a new Queue
func NewQueue() *Queue {
	return &Queue{
		messages: make([]string, 0),
	}
}

// Enqueue adds message to the queue
func (q *Queue) Enqueue(msg string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.messages = append(q.messages, msg)
}

// Dequeue remove and return first element in queue
func (q *Queue) Dequeue() (string, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.messages) == 0 {
		return "", false
	}

	msg := q.messages[0]
	q.messages = q.messages[1:]

	return msg, true
}

// Length return the number of message in the queue
func (q *Queue) Length() int {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	return len(q.messages)
}
