package plog

import (
	"github.com/crowdflux/angel/app/plog/log_tags"
	"tag"
)

type message struct {
	Tag    tag.Tag     `json:"Tag"`
	Params interface{} `json:"params,omitempty"`
}

func MessageWithParam(tag *tag.Tag, param interface{}) message {
	return message{
		Tag:    *tag,
		Params: param,
	}
}

func Message(param interface{}) message {
	return message{
		Tag:    *log_tags.MESSAGE,
		Params: param,
	}
}
