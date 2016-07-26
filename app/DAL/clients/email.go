package clients

import (
	"github.com/crowdflux/angel/app/config"
	"github.com/sendgrid/sendgrid-go"
)

// Donot use this for now
func NewSendGridClient() *sendgrid.SGClient {
	client := sendgrid.NewSendGridClientWithApiKey(config.Get(config.SENDGRID_API_KEY))
	return client
}
