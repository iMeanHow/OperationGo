package main

import "HelloGo/OperationGo/RabbitMQ"

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQDirect("DirectQueue")
	rabbitmq.ConsumeDirect()
}
