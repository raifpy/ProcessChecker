package ofisclient

import "context"

type Config struct {
	DialAddr string
	Context  context.Context `json:"-"`
}

type WebsocketResponse struct {
	MessageID string
	Message   interface{}
}

type WebsocketRequest struct {
	ReplyID string
	Message interface{}
}

type ResponseMessage map[string]interface{}

/*
type ResponseMessageProcessCheck struct {
	ProcessName string
}

type ResponseMessageError struct {
	Error string
}
*/

type Sorgu struct {
	Hostname string
	Process  bool
	Error    string
}
