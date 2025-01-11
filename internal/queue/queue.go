package queue

import (
	"simple-rabbit/internal/entities"
	"simple-rabbit/internal/storage"
	"sync"
	"time"
)

type Queue struct {
	messages []entities.Message
	mutex    sync.Mutex
	storage  *storage.Storage
	name     string
}

// NewQueue create a new Queue
func NewQueue(name string, storage *storage.Storage) *Queue {
	q := &Queue{
		messages: make([]entities.Message, 0),
		name:     name,
		storage:  storage,
	}

	if storage != nil {
		storedMessage, err := storage.LoadMessages(name)
		if err != nil {
			panic(err)
		}
		q.messages = append(q.messages, storedMessage...)
	}

	return q
}

// Enqueue adds message to the queue
func (q *Queue) Enqueue(msg entities.Message) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	//expiration := time.Now().Add(ttl)

	q.messages = append(q.messages, msg)
	if q.storage != nil {
		_ = q.storage.SaveMessage(q.name, msg)
	}
}

// Dequeue remove and return first element in queue
func (q *Queue) Dequeue() (*entities.Message, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.messages) == 0 {
		return nil, false
	}

	for i, msg := range q.messages {
		// Check if the message has expired
		if msg.Expiration != nil && msg.Expiration.Before(time.Now()) {
			// Skip expired messages
			if msg.Key != nil {
				q.storage.DeleteMessage(q.name, msg.Key)
			}

			continue
		}

		// Remove the message from the queue
		q.messages = append(q.messages[:i], q.messages[i+1:]...)

		// Delete the message from storage if it has a key
		if msg.Key != nil {
			q.storage.DeleteMessage(q.name, msg.Key)
		}

		// Return the valid message
		return &msg, true

	}

	return nil, false
}

// Length return the number of message in the queue
func (q *Queue) Length() int {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	return len(q.messages)
}
