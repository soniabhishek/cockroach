package feed_line

import (
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/models"
	"github.com/streadway/amqp"
	"sync"
)

// Feedline logger channel
type feedlineLoggerChannel struct {
	mq rabbitmq.MQ

	queueName string
	once      sync.Once

	// purpose of this is for confirming logs have been written to db
	latestMsg amqp.Delivery
}

func newFeedLineLogger() *feedlineLoggerChannel {

	name := "feedlinelogQueue"

	mq := rabbitmq.New(name)

	go func() {

		for flog := range tempFllChan {

			// Send only the models.Feedline part of the flu in bytes
			bty, _ := json.Marshal(flog)
			mq.Publish(bty)

			//plog.Info("feedline logger", "complete push from: ", fll.queueName, "id: ", flog.FluId.String())
		}
	}()

	return &feedlineLoggerChannel{
		queueName: name,
		mq:        mq,
	}
}

// Ask purpose of this to himanshu
var tempFllChan = make(chan models.FeedLineLog)

func (fll *feedlineLoggerChannel) Push(flog models.FeedLineLog) {

	tempFllChan <- flog

}

func (fll *feedlineLoggerChannel) Receiver() <-chan models.FeedLineLog {

	println("Feedline, subscribe request: ", fll.queueName)

	var flogChan chan models.FeedLineLog
	var flag bool = false

	fll.once.Do(func() {

		flogChan = make(chan models.FeedLineLog)

		go func() {

			for msg := range fll.mq.Consume() {

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
