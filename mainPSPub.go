package main

import (
	"HelloGo/OperationGo/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" + "newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("P/S published No." + strconv.Itoa(i) + " message")
		fmt.Println("P/S produced No." + strconv.Itoa(i) + " message")
		time.Sleep(1 * time.Second)
	}
}