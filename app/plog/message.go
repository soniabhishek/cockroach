package plog

type Message struct {
	// Required
	Message string `json:"message"`

	// Optional
	Params []interface{} `json:"params,omitempty"`
}

func (m *Message) Class() string { return "logentry" }
