package main

import (
	"fmt"
	"simple-rabbit/internal/broker"
	"simple-rabbit/internal/entities"
)

func main() {
	b := broker.NewBroker("broker2.db")

	// create queue
	b.CreateQueue("order")

	//create exchange
	b.CreateExchange("main_exchange")

	//bind queue
	b.BindQueue("main_exchange", "order.created", "order")

	//send message
	b.SendMessage("main_exchange", "order.created", entities.Message{Content: "new order created !"})

	b.ReceiveMessage("order")

	// verify persist message
	b = broker.NewBroker("broker.db")
	b.ReceiveMessage("main_exchange")

	message := b.ReceiveMessage("order")
	fmt.Printf("Receive message from queue order: %s\n", message)
}
