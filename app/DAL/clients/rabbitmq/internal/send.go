package rabbitmq

//
//package main
//
//import (
//	"log"
//
//	"encoding/json"
//	"github.com/crowdflux/angel/app/models"
//	"github.com/crowdflux/angel/app/models/uuid"
//	"github.com/streadway/amqp"
//)
//
//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Fatalf("%s: %s", msg, err)
//	}
//}
//
//func main() {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	failOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//	q, err := ch.QueueDeclare(
//		"hello", // name
//		false,   // durable
//		false,   // delete when unused
//		false,   // exclusive
//		false,   // no-wait
//		nil,     // arguments
//	)
//	failOnError(err, "Failed to declare a queue")
//
//	flu := models.FeedLineUnit{
//		ID: uuid.NewV4(),
//	}
//	fluBty, _ := json.Marshal(flu)
//
//	err = ch.Publish(
//		"",     // exchange
//		q.Name, // routing key
//		false,  // mandatory
//		false,  // immediate
//		amqp.Publishing{
//			ContentType: "application/json",
//			Body:        fluBty,
//		})
//	log.Printf(" [x] Sent %s", flu)
//	failOnError(err, "Failed to publish a message")
//
//}
