package rabbitmq

import (
	//"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBhosda_Publish(t *testing.T) {

	test := New("test")

	bty := []byte("whore")

	test.Publish(bty)

	test.Consume()

	time.Sleep(time.Duration(1) * time.Minute)

	//assert.Equal(t, bty, msg)

}
