package broker

import (
	"fmt"
	"simple-rabbit/internal/exchange"
	"simple-rabbit/internal/queue"
	"simple-rabbit/internal/storage"
)

type Broker struct {
	queues    map[string]*queue.Queue
	exchanges map[string]*exchange.Exchange
	storage   *storage.Storage
}

func NewBroker(storagePath string) *Broker {
	return &Broker{
		queues:    make(map[string]*queue.Queue),
		exchanges: make(map[string]*exchange.Exchange),
		storage:   storage.NewStorage(storagePath),
	}
}

// CreateQueue create queue
func (b *Broker) CreateQueue(name string) {
	if _, ok := b.queues[name]; ok {
		fmt.Printf("Queue %s already exists", name)
		return
	}

	b.queues[name] = queue.NewQueue(name, b.storage)
	fmt.Printf("Queue %s created", name)
}

// SendMessage send message to queue
func (b *Broker) SendMessage(exChange, routingKey, message string) {
	ex, ok := b.exchanges[exChange]
	if !ok {
		fmt.Printf("Exchange %s not found", exChange)
		return
	}

	ex.PublishMessage(routingKey, message)

}

// ReceiveMessage receive message
func (b *Broker) ReceiveMessage(queueName string) any {
	if q, ok := b.queues[queueName]; ok {
		if message, ok := q.Dequeue(); ok {
			fmt.Printf("Message Received from queue:%s message:%s", queueName, message)
			return message
		} else {
			fmt.Printf("Queue %s not found", queueName)
		}
	} else {
		fmt.Printf("Queue %s not found", queueName)
	}

	return nil
}

// CreateExchange create exchange
func (b *Broker) CreateExchange(name string) {
	if _, ok := b.exchanges[name]; ok {
		fmt.Printf("Exchange %s already exists", name)
		return
	}

	b.exchanges[name] = exchange.NewExchange(name)
	fmt.Printf("Exchange %s created", name)
}

// BindQueue bind queue to routing key
func (b *Broker) BindQueue(exchange, routingKey, queueName string) {
	ex, exist := b.exchanges[exchange]
	if !exist {
		fmt.Printf("Exchange %s not found", exchange)
		return
	}

	q, exist := b.queues[queueName]
	if !exist {
		fmt.Printf("Queue %s not found", queueName)
		return
	}

	ex.BindToQueue(routingKey, q)
}
