package main

import (
	"fmt"
	"simple-rabbit/internal/broker"
)

func main() {
	b := broker.NewBroker("broker.db")

	// create queue
	b.CreateQueue("order")

	//create exchange
	b.CreateExchange("main_exchange")

	//bind queue
	b.BindQueue("main_exchange", "order.created", "order")

	//send message
	b.SendMessage("main_exchange", "order.created", "new order created !")

	// verify persist message
	b = broker.NewBroker("broker.db")
	b.ReceiveMessage("main_exchange")

	message := b.ReceiveMessage("order")
	fmt.Printf("Receive message from queue order: %s\n", message)
}
