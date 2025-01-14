package exchange

import (
	"fmt"
	"simple-rabbit/internal/entities"
	"simple-rabbit/internal/queue"
)

type Exchange struct {
	Name       string
	QueueMap   map[string]*queue.Queue
	QueueNames []string
}

func NewExchange(name string) *Exchange {
	return &Exchange{
		Name:     name,
		QueueMap: make(map[string]*queue.Queue),
	}
}

// BindToQueue bind a queue to specific routing key
func (e *Exchange) BindToQueue(routingKey string, q *queue.Queue) {
	e.QueueMap[routingKey] = q
	e.QueueNames = append(e.QueueNames, routingKey)
	fmt.Printf("queue bounde to Exchange %s with routing key  %s\n", e.Name, routingKey)
}

// PublishMessage publish a message with a routing key to the bounded queue
func (e *Exchange) PublishMessage(routingKey string, message entities.Message) {
	if q, exists := e.QueueMap[routingKey]; exists {
		q.Enqueue(message)
		fmt.Printf("queue bounde to Exchange %s with routing key  %s\n", e.Name, routingKey)
	} else {
		fmt.Printf("no queue bounde to routing key: %s ", routingKey)
	}
}
