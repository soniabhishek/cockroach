package rabbitmq

import (
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/plog"
	"github.com/streadway/amqp"
	"log"
)

// message is the application type for a message.  This can contain identity,
// or a reference to the recevier chan for further demuxing.
type message []byte

// session composes an amqp.Connection with an amqp.Channel
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// redial continually connects to the URL, exiting the program when no longer possible
func redial() chan chan session {
	sessions := make(chan chan session)

	username := config.RABBITMQ_USERNAME.Get()
	password := config.RABBITMQ_PASSWORD.Get()
	host := config.RABBITMQ_HOST.Get()

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {
			sessions <- sess

			conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":5672/")
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, username)
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}

			sess <- session{conn, ch}
		}
	}()

	return sessions
}

// publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func publish(sessions chan chan session, messages <-chan message, qName string) {
	var (
		pending = make(chan message, 1)
	)

	for session := range sessions {
		pub := <-session

		plog.Trace("publishing...")

		if _, err := pub.QueueDeclare(qName, true, false, false, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", qName, err)
			return
		}

	loop:
		for {
			select {
			case msg := <-messages:
				err := pub.Publish(
					"",    // exchange
					qName, // routing key
					false, // mandatory
					false, // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        msg,
					})
				if err != nil {

					log.Println("dscs")
					pending <- msg
					pub.Close()
					break loop
				}
			case msg := <-pending:
				err := pub.Publish(
					"",    // exchange
					qName, // routing key
					false, // mandatory
					false, // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        msg,
					})
				if err != nil {
					log.Println("dscs")

					pending <- msg
					pub.Close()
					break loop
				}
			}
		}
	}
}

// subscribe consumes deliveries from an exclusive queue from a fanout exchange and sends to the application specific messages chan.
func subscribe(sessions chan chan session, messages chan<- amqp.Delivery, qName string) {

	for session := range sessions {
		sub := <-session

		deliveries, err := sub.Consume(qName, "", false, false, false, false, nil)
		if err != nil {
			log.Printf("cannot consume from: %q, %v", qName, err)
			return
		}

		plog.Trace("subscribed...")

		for msg := range deliveries {
			messages <- msg
		}
	}
}
