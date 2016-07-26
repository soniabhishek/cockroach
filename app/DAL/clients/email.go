package clients

import (
	"github.com/crowdflux/angel/app/config"
	"gopkg.in/sendgrid/sendgrid-go.v2"
)

// Donot use this for now
func NewSendGridClient() *sendgrid.SGClient {
	client := sendgrid.NewSendGridClientWithApiKey(config.SENDGRID_API_KEY.Get())
	return client
}
