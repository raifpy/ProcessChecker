package websocket

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func (c Client) Write(message ClientRequest) error {
	if c.Id == "" {
		return errors.New("id cannot be empty")
	}

	return c.Conn.WriteJSON(message)
}

func (c Client) WriteWithID(m ClientRequest) (id string, err error) {
	id = RandStringRunes(10) + "_messageid"
	m.MessageID = id
	err = c.Write(m)
	return
}

func (ws *Ws) Write(c Client, m ClientRequest) (id string, err error) {
	if m.MessageID == "" {
		m.MessageID = RandStringRunes(10) + "_wsmessage"
	}
	id = m.MessageID

	fmt.Println("write called")

	ws.Redirector.SetAsTimeout(m.MessageID, make(chan ClientParsedResponse), ws.Redirector.Timeout)
	if err := c.Write(m); err != nil {
		ws.Redirector.Del(m.MessageID)
		return "", err
	}

	return
}

func (ws *Ws) WriteAndRead(c Client, m ClientRequest) (chan ClientParsedResponse, error) {
	id, err := ws.Write(c, m)
	if err != nil {
		return nil, err
	}

	ch, ok := ws.Redirector.Get(id)
	if !ok {
		return nil, errors.New("redirector id not exists")
	}
	return ch, nil
}

func (ws *Ws) WriteAndWait(c Client, m ClientRequest) (ClientParsedResponse, error) {
	id, err := ws.Write(c, m)
	if err != nil {
		return ClientParsedResponse{}, err
	}

	ch, ok := ws.Redirector.Get(id)
	if !ok {
		return ClientParsedResponse{}, errors.New("redirector id not exists")
	}
	defer func() {
		ws.Redirector.Del(id)
		close(ch)
	}()
	select {
	case <-time.NewTicker(ws.Redirector.Timeout).C:
		{
			return ClientParsedResponse{}, errors.New("timeout")
		}
	case v := <-ch:
		{
			return v, nil
		}
	}
}

func (ws *Ws) WriteAll(m ClientRequest) []ClientParsedResponse {
	list := []ClientParsedResponse{}
	ws.Clients.mutex.RLock()
	var wait sync.WaitGroup
	for id, client := range ws.Clients.clientmaps {
		wait.Add(1)
		go func(id string, client Client) {
			resp, err2 := ws.WriteAndWait(client, m)
			list = append(list, ClientParsedResponse{
				Id:      id,
				Time:    time.Now(),
				Message: resp,
				Error:   err2,
			})

			wait.Done()

		}(id, client)
	}
	ws.Clients.mutex.RUnlock()

	wait.Wait()

	return list
}
