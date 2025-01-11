package queue

import (
	"container/heap"
	"sync"
	"time"

	"simple-rabbit/internal/entities"
	"simple-rabbit/internal/storage"
)

type Queue struct {
	pq      *PriorityQueue
	mutex   sync.Mutex
	storage *storage.Storage
	name    string
}

// PriorityQueue implements heap.Interface
type PriorityQueue struct {
	messages []entities.Message
}

func (pq PriorityQueue) Len() int {
	return len(pq.messages)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq.messages[i].Priority > pq.messages[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.messages[i], pq.messages[j] = pq.messages[j], pq.messages[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(entities.Message)
	pq.messages = append(pq.messages, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.messages
	n := len(old)
	item := old[n-1]
	pq.messages = old[0 : n-1]
	return item
}

// NewQueue create a new Queue
func NewQueue(name string, storage *storage.Storage) *Queue {
	q := &Queue{
		pq:      &PriorityQueue{messages: make([]entities.Message, 1000)}, //1000  for prevent regenerate slice
		name:    name,
		storage: storage,
	}

	if storage != nil {
		storedMessages, err := storage.LoadMessages(name)
		if err != nil {
			panic(err)
		}

		now := time.Now()
		validMessages := make([]entities.Message, 0, len(storedMessages))
		for _, msg := range storedMessages {
			if msg.Expiration == nil || msg.Expiration.After(now) {
				validMessages = append(validMessages, msg)
			} else if msg.Key != nil {
				storage.DeleteMessage(name, msg.Key)
			}
		}

		q.pq.messages = validMessages
		heap.Init(q.pq)
	}

	return q
}

// Enqueue adds message to the queue
func (q *Queue) Enqueue(msg entities.Message) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	heap.Push(q.pq, msg)

	//expiration := time.Now().Add(ttl)
	if q.storage != nil {
		_ = q.storage.SaveMessage(q.name, msg)
	}
}

// Dequeue remove and return first element in queue
func (q *Queue) Dequeue() (*entities.Message, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.pq.messages) == 0 {
		return nil, false
	}

	// پیدا کردن اولین پیام غیر منقضی
	for q.pq.Len() > 0 {
		msg := heap.Pop(q.pq).(entities.Message)

		// بررسی انقضا
		if msg.Expiration != nil && msg.Expiration.Before(time.Now()) {
			if msg.Key != nil {
				q.storage.DeleteMessage(q.name, msg.Key)
			}
			continue
		}

		// حذف از storage
		if msg.Key != nil {
			q.storage.DeleteMessage(q.name, msg.Key)
		}

		return &msg, true
	}

	for i, msg := range q.pq.messages {
		// Check if the message has expired
		if msg.Expiration != nil && msg.Expiration.Before(time.Now()) {
			// Skip expired pq
			if msg.Key != nil {
				q.storage.DeleteMessage(q.name, msg.Key)
			}

			continue
		}

		// Remove the message from the queue
		q.pq.messages = append(q.pq.messages[:i], q.pq.messages[i+1:]...)

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

	return len(q.pq.messages)
}
