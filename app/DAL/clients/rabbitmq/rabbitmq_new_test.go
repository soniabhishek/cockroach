package rabbitmq

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func TestBhosda_Publish(t *testing.T) {

	test := New("test")

	go func() {
		i := 0
		for {
			bty := []byte("cho" + strconv.Itoa(i))
			i++
			test.Publish(bty)
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()

	go func() {
		for msg := range test.Consume() {
			log.Println(string(msg.Body))
			msg.Ack(false)
		}
	}()

	//assert.EqualValues(t, bty, msg)

	time.Sleep(time.Duration(1) * time.Minute)

}
