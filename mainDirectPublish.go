package main

import (
	"HelloGo/OperationGo/RabbitMQ"
	"fmt"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQDirect("DirectQueue")
	rabbitmq.PublishDirect("Direct Exchange Test")
	fmt.Println("Message sent!")
}
