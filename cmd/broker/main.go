package main

import (
	"simple-rabbit/internal/broker"
	"simple-rabbit/internal/entities"
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
	b.SendMessage("main_exchange", "order.created", entities.Message{Content: "new order(5) created !", Priority: 5})
	b.SendMessage("main_exchange", "order.created", entities.Message{Content: "new order(2) created !", Priority: 2})
	b.SendMessage("main_exchange", "order.created", entities.Message{Content: "new order(4) created !", Priority: 4})
	b.SendMessage("main_exchange", "order.created", entities.Message{Content: "new order(3) created !", Priority: 3})
	b.SendMessage("main_exchange", "order.created", entities.Message{Content: "new order(1) created !", Priority: 1})

	b.ReceiveMessage("order")
	b.ReceiveMessage("order")
	b.ReceiveMessage("order")
	b.ReceiveMessage("order")
	b.ReceiveMessage("order")

}
