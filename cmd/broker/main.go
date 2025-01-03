package broker

import "simple-rabbit/internal/broker"

func main() {
	b := broker.NewBroker()

	b.CreateQueue("test-queue")

	b.SendMessage("test-queue", "hello world")
	b.SendMessage("test-queue", "this is a test message")

	b.ReceiveMessage("test-queue")
	b.ReceiveMessage("test-queue")
}
