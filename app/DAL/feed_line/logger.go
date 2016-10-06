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

// Feedline logger channel
type feedlineLoggerChannel struct {
	amqpChan  *amqp.Channel
	queueName string
	once      sync.Once

	// purpose of this is for confirming logs have been written to db
	latestMsg amqp.Delivery
}

func newFeedLineLogger() *feedlineLoggerChannel {

	name := "feedlinelogQueue"

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

	return &feedlineLoggerChannel{
		amqpChan:  ch,
		queueName: q.Name,
	}
}

func (fll *feedlineLoggerChannel) Push(flog models.FeedLineLog) {

	// Send only the models.Feedline part of the flu in bytes
	bty, _ := json.Marshal(flog)

	// This is async
	// TODO Think about a way to guarantee this operation also
	err := fll.amqpChan.Publish(
		"",            // exchange
		fll.queueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bty,
		})
	if err != nil {
		plog.Error("Feedline logger", err, "error publishing to channel", "flu_id: "+flog.FluId.String())
		panic(err)
	}

	//plog.Info("feedline logger", "complete push from: ", fll.queueName, "id: ", flog.FluId.String())
}

func (fll *feedlineLoggerChannel) Receiver() <-chan models.FeedLineLog {

	println("Feedline, subscribe request: ", fll.queueName)

	var flogChan chan models.FeedLineLog
	var flag bool = false

	fll.once.Do(func() {

		flogChan = make(chan models.FeedLineLog)

		deliveryChan, err := fll.amqpChan.Consume(
			fll.queueName, // queue
			"",            // consumer
			false,         // auto-ack
			false,         // exclusive
			false,         // no-local
			false,         // no-wait
			nil,           // args
		)
		if err != nil {
			plog.Error("Feedline", err, "error consuming queue, name:", fll.queueName)
			panic(err)
		}

		go func() {

			for msg := range deliveryChan {

				flog := models.FeedLineLog{}
				json.Unmarshal(msg.Body, &flog)

				fll.latestMsg = msg

				flogChan <- flog
				//plog.Info("feedline logger", "sent to FLU chan, name: ", fll.queueName, "id: ", flog.FluId.String())
			}
		}()

		flag = true
	})

	if flag {
		return (<-chan models.FeedLineLog)(flogChan)
	} else {
		panic(errors.New("Feedline logger already subscribed, name: " + fll.queueName))
	}
}

func (fll *feedlineLoggerChannel) ConfirmReceivedProcessed() {
	err := fll.latestMsg.Ack(true)
	if err != nil {
		panic(err)
	}
}

var stdLoggerChan *feedlineLoggerChannel

var once sync.Once

func GetFeedlineLoggerChannel() *feedlineLoggerChannel {

	once.Do(func() {
		stdLoggerChan = newFeedLineLogger()
	})

	return stdLoggerChan
}
