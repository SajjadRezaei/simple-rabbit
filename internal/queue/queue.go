package queue

import (
	"simple-rabbit/internal/storage"
	"sync"
)

type Queue struct {
	messages []string
	mutex    sync.Mutex
	storage  *storage.Storage
	name     string
}

// NewQueue create a new Queue
func NewQueue(name string, storage *storage.Storage) *Queue {
	q := &Queue{
		messages: make([]string, 0),
		name:     name,
		storage:  storage,
	}

	if storage == nil {
		storedMessage, _ := storage.LoadMessage(name)
		q.messages = append(q.messages, storedMessage...)
	}

	return q
}

// Enqueue adds message to the queue
func (q *Queue) Enqueue(msg string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = append(q.messages, msg)
	if q.storage == nil {
		_ = q.storage.SaveMessage(q.name, msg)
	}
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
