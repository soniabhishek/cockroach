package rabbitmq

import (
	"github.com/streadway/amqp"
)

type MQ struct {
	name     string
	messages chan message
}

func (b *MQ) Publish(bty []byte) {
	b.messages <- message(bty)
}

func (b *MQ) Consume() <-chan amqp.Delivery {

	messages := make(chan amqp.Delivery)
	go subscribe(redial(), messages, b.name)
	return messages
}

func New(name string) MQ {
	b := MQ{name, make(chan message)}

	go publish(redial(), b.messages, b.name)

	return b
}
