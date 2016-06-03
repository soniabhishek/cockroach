package clients

import (
	"github.com/sendgrid/sendgrid-go"
	"gitlab.com/playment-main/angel/app/config"
)

// Donot use this for now
func NewSendGridClient() *sendgrid.SGClient {
	client := sendgrid.NewSendGridClientWithApiKey(config.Get(config.SENDGRID_API_KEY))
	return client
}
