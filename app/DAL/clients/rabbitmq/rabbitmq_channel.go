package rabbitmq

import (
	"github.com/crowdflux/angel/app/config"
	"github.com/streadway/amqp"
)

var rabbitmqConn *amqp.Connection

func init() {
	rabbitmqConn = initRabbitMqClient()
}

func initRabbitMqClient() *amqp.Connection {

	username := config.RABBITMQ_USERNAME.Get()
	password := config.RABBITMQ_PASSWORD.Get()
	host := config.RABBITMQ_HOST.Get()

	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":5672/")
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
