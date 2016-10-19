package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"github.com/crowdflux/angel/app/config"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"time"
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
func redial(exchange string, ctx context.Context) chan chan session {
	sessions := make(chan chan session)

	username := config.RABBITMQ_USERNAME.Get()
	password := config.RABBITMQ_PASSWORD.Get()
	host := config.RABBITMQ_HOST.Get()

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {

			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Println("shutting down session factory")
				return
			}

			conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":5672/")
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, username)
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}

			if err := ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil); err != nil {
				log.Fatalf("cannot declare fanout exchange: %v", err)
			}

			select {
			case sess <- session{conn, ch}:
			case <-ctx.Done():
				log.Println("shutting down new session")
				return
			}
		}
	}()

	return sessions
}

// publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func publish(sessions chan chan session, messages <-chan message, exchange string) {
	var (
		running bool
		reading = messages
		pending = make(chan message, 1)
		confirm = make(chan amqp.Confirmation, 1)
	)

	for session := range sessions {
		pub := <-session

		// publisher confirms for this channel/connection
		if err := pub.Confirm(false); err != nil {
			log.Printf("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			pub.NotifyPublish(confirm)
		}

		log.Printf("publishing...")

		for {
			var body message
			select {
			case confirmed := <-confirm:
				log.Println("1")
				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, string(body))
				}
				reading = messages

			case body = <-pending:
				log.Println("2")
				routingKey := "ignored for fanout exchanges, application dependent for other exchanges"
				err := pub.Publish(exchange, routingKey, false, false, amqp.Publishing{
					Body: body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					pub.Close()
					break
				}

			case body, running = <-reading:
				log.Println("3")

				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}

// identity returns the same host/process unique string for the lifetime of
// this process so that subscriber reconnections reuse the same queue name.
func identity() string {
	hostname, err := os.Hostname()
	h := sha1.New()
	fmt.Fprint(h, hostname)
	fmt.Fprint(h, err)
	fmt.Fprint(h, os.Getpid())
	return fmt.Sprintf("%x", h.Sum(nil))
}

// subscribe consumes deliveries from an exclusive queue from a fanout exchange and sends to the application specific messages chan.
func subscribe(sessions chan chan session, messages chan<- message, qName, exchange string) {
	queue := identity()

	for session := range sessions {
		sub := <-session

		if _, err := sub.QueueDeclare(queue, true, false, false, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", queue, err)
			return
		}

		routingKey := "application specific routing key for fancy toplogies"
		if err := sub.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
			log.Printf("cannot consume without a binding to exchange: %q, %v", exchange, err)
			return
		}

		deliveries, err := sub.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			log.Printf("cannot consume from: %q, %v", queue, err)
			return
		}

		log.Printf("subscribed...")

		for msg := range deliveries {
			log.Println("1'")
			messages <- message(msg.Body)
			log.Println("2'")
			sub.Ack(msg.DeliveryTag, false)
		}
	}
}

// read is this application's translation to the message format, scanning from
// stdin.
func read(r io.Reader) <-chan message {
	lines := make(chan message)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- message(scan.Bytes())
		}
	}()
	return lines
}

// write is this application's subscriber of application messages, printing to
// stdout.
func write(w io.Writer) chan<- message {
	lines := make(chan message)
	go func() {
		for line := range lines {
			fmt.Fprintln(w, string(line))
		}
	}()
	return lines
}

type Bhosda struct {
	name     string
	messages chan message
	ctx      context.Context
}

func (b *Bhosda) Publish(bty []byte) {
	b.messages <- message(bty)
}

func (b *Bhosda) Consume() <-chan message {

	messages := make(chan message)
	go subscribe(redial(b.name, b.ctx), write(os.Stdout), b.name, b.name)
	return messages
}

func New(name string) Bhosda {
	ctx, _ := context.WithCancel(context.Background())
	b := Bhosda{name, make(chan message), ctx}

	go publish(redial(b.name, b.ctx), b.messages, b.name)

	return b
}

func main() {

	b := New("test123")

	b.Consume()

	time.Sleep(time.Duration(3) * time.Second)

	b.Publish([]byte("sdvsdv1"))

	time.Sleep(time.Duration(10) * time.Minute)
}
