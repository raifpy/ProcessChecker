package websocket

import (
	"encoding/json"
	"time"
)

type ClientParsedResponse struct {
	Id      string
	Time    time.Time
	Message interface{}
	Error   error
}

type ClientResponse struct {
	ReplyID string
	Message interface{}
}

type ClientRequest struct {
	MessageID string
	Message   interface{}
}

func (cr ClientRequest) ToJson() []byte {
	a, _ := json.Marshal(cr)
	return a
}

type MessageError struct {
	Error string
}

type RequestMessage struct {
	RequestMessageError
	RequestMessageProcessCheck
}

type RequestMessageProcessCheck struct {
	ProcessName string
}

type RequestMessageError struct {
	Error string
}

type Sorgu struct {
	Hostname string
	Process  bool
	Error    string
}
