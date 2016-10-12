package feed_line

import (
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/streadway/amqp"
	"sync"
)

// ShortHand for channel of FLUs i.e. FeedLine
type Fl struct {
	amqpChan  *amqp.Channel
	queueName string
	once      sync.Once
}

func New(name string) Fl {

	ch := rabbitmq.GetNewChannel()

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		plog.Error("Feedline", err, "error declaring queue, name: ", name)
		panic(err)
	}

	return Fl{
		amqpChan:  ch,
		queueName: q.Name,
	}
}

func (fl *Fl) Push(flu FLU) {

	// Send only the models.Feedline part of the flu in bytes
	bty, _ := json.Marshal(flu.FeedLineUnit)

	// This is async
	// TODO Think about a way to guarantee this operation also
	err := fl.amqpChan.Publish(
		"",           // exchange
		fl.queueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bty,
		})
	if err != nil {
		plog.Error("Feedline", err, "error publishing to channel", "flu_id: "+flu.ID.String())
		panic(err)
	}

	// Just for safety: if someone forgets
	// to ConfirmReceive the flu received from a queue
	// then reconfirm it here as it will most
	// probably be a bug
	if flu.delivery.Acknowledger != nil {
		flu.ConfirmReceive()
	}

	plog.Info("feedline", "complete push from: ", fl.queueName, "id: ", flu.ID.String())
}

func (fl *Fl) Receiver() <-chan FLU {

	println("Feedline, subscribe request: ", fl.queueName)

	var fluChan chan FLU
	var flag bool = false

	fl.once.Do(func() {

		fluChan = make(chan FLU)

		deliveryChan, err := fl.amqpChan.Consume(
			fl.queueName, // queue
			"",           // consumer
			false,        // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			plog.Error("Feedline", err, "error consuming queue, name:", fl.queueName)
			panic(err)
		}

		go func() {

			for msg := range deliveryChan {

				flu := models.FeedLineUnit{}
				json.Unmarshal(msg.Body, &flu)

				fluChan <- FLU{
					FeedLineUnit: flu,
					delivery:     msg,
					once:         &sync.Once{},
				}
				plog.Info("feedline", "sent to FLU chan, name: ", fl.queueName, "id: ", flu.ID.String())
			}
		}()

		flag = true
	})

	if flag {
		return (<-chan FLU)(fluChan)
	} else {
		panic(errors.New("Feedline already subscribed, name: " + fl.queueName))
	}

}
