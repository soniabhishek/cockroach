package feed_line

import (
	"encoding/json"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
)

type feedLineLogChan struct {
}

func newLoggerChan() feedLineLogChan {

	ch := rabbitmq.GetNewChannel()

	q, err := ch.QueueDeclare(
		"FeedlineLogQueue", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	if err != nil {
		plog.Error("Feedline Logger", err, "error declaring queue")
		panic(err)
	}

	outChan, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		plog.Error("Feedline", err, "error consuming queue")
	}

	fl := Fl{make(chan FLU, 50000), outChan, ch, q.Name, false}

	go func() {

		for msg := range outChan {

			flu := models.FeedLineUnit{}
			json.Unmarshal(msg.Body, &flu)

			fl.ch <- FLU{
				FeedLineUnit: flu,
				delivery:     msg,
			}
		}

	}()

	return fl
}
