package main

import "HelloGo/OperationGo/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" + "newProduct")
	rabbitmq.RecieveSub()
}
