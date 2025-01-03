package broker

import (
	"fmt"
	"simple-rabbit/internal/queue"
)

type Broker struct {
	queues map[string]*queue.Queue
}

func NewBroker() *Broker {
	return &Broker{
		queues: make(map[string]*queue.Queue),
	}
}

func (b *Broker) CreateQueue(name string) {
	if _, ok := b.queues[name]; ok {
		fmt.Printf("Queue %s already exists", name)
		return
	}

	b.queues[name] = queue.NewQueue()
	fmt.Printf("Queue %s created", name)
}

func (b *Broker) SendMessage(queueName, message string) {
	if q, ok := b.queues[queueName]; ok {
		q.Enqueue(message)
		fmt.Printf("Message sent to queue %s", queueName)
		return
	}

	fmt.Printf("Queue %s not found", queueName)

}

func (b *Broker) ReceiveMessage(queueName string) {
	if q, ok := b.queues[queueName]; ok {
		if message, ok := q.Dequeue(); ok {
			fmt.Printf("Message Received from queue:%s message:%s", queueName, message)
		} else {
			fmt.Printf("Queue %s not found", queueName)
		}
	} else {
		fmt.Printf("Queue %s not found", queueName)
	}
}
