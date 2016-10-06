package rabbitmq

import "github.com/streadway/amqp"

var rabbitmqConn *amqp.Connection

func init() {
	rabbitmqConn = initRabbitMqClient()
}

func initRabbitMqClient() *amqp.Connection {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	return conn
}

func GetNewChannel() *amqp.Channel {
	ch, err := rabbitmqConn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
