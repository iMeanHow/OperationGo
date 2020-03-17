package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//Connection for RabbitMQ
const RMQurl  = "amqp://iMeanHow1:simpletest@127.0.0.1:5672/Test1" //amqp:Username:Password@IP:Port/VirtualHost

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	//Queue Name
	QueueName string
	//Exchange
	Exchange string
	//Key
	Key string
	//Connection info
	Mqurl string
}

// Instantiate RabbitMQ
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ{
	return &RabbitMQ{QueueName:queueName, Exchange:exchange, Key:key, Mqurl:RMQurl}


}

//Disable RabbitMQ Connection
func (r *RabbitMQ) DisableConn(){
	r.channel.Close()
	r.conn.Close()
}

//Error Handling
func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Fatal("%s:%s", message, err) //Fatal = print + exit(1). No defer() performed
		panic(fmt.Sprintf("%s:%s", message, err)) //Similar to throw, defer() will be performed
	}
}

//Instantiate Direct Exchange mode
func NewRabbitMQDirect(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "") //"" equals default setting
	//Build Connection with RabbitMQ
	var err error
	rabbitmq.conn,err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "New connection failed!")
	rabbitmq.channel,err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "Obtain channel failed!")
	return rabbitmq

}

func (r *RabbitMQ) PublishDirect(message string) {
	//1. Apply for the queue
	//   Exist -> Apply
	//   Not Exist -> Create -> Apply
	// QueueName, Durable, autoDelete, exclusive, noWait, args
	_,err := r.channel.QueueDeclare(r.QueueName,false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}

	//2. Send message to queue
	r.channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType:"text/plain",
		Body:[]byte(message),
	})
}

func (r *RabbitMQ) ConsumeDirect() {
	_,err := r.channel.QueueDeclare(r.QueueName,false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}

	msgs,err := r.channel.Consume(r.QueueName, "", true, false, false, false, nil,)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//Goroutine Process Msgs
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messasges, to exit press Ctrl+^C")
	<- forever
}


//Instantiate Fanout Exchange Mode (Pub/Sub)
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "") //"" equals default setting
	//Build Connection with RabbitMQ
	var err error
	rabbitmq.conn,err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "New connection failed!")
	rabbitmq.channel,err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "Obtain channel failed!")
	return rabbitmq
}

//PubSub Exchange Mode (Fannout)
func (r *RabbitMQ) PublishPub(message string) {
	err := r.channel.ExchangeDeclare(r.Exchange, "fanout", true, false, false, false, nil)

	r.failOnError(err, "P/S Exchange declare failed")

	err = r.channel.Publish(r.Exchange, "", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//Sub-Consume
func (r *RabbitMQ) RecieveSub() {
	err := r.channel.ExchangeDeclare(r.Exchange, "fanout", true, false, false, false, nil)
	r.failOnError(err, "P/S Exchange declare failed")

	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnError(err, "Failed to declare a queue")

	err = r.channel.QueueBind(q.Name, "", r.Exchange, false, nil)

	messges, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[*] Waiting for messasges, to exit press Ctrl+^C")
	<-forever
}