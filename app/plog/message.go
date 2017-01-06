package plog

type message struct {
	// Required
	Message string `json:"message"`

	// Optional
	Params []interface{} `json:"params,omitempty"`
}

func (m message) Class() string { return "logentry" }

func NewMessageWithParam(msg string, param interface{}) message {
	return message{
		Message: msg,
		Params:  []interface{}{param},
	}
}

func NewMessage(msg string) message {
	return message{
		Message: msg,
	}
}
