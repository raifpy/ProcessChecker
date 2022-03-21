package websocket

import (
	"errors"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
}

type WsClients struct {
	clientmaps map[string]Client
	mutex      *sync.RWMutex
}

func NewWsClients() *WsClients {
	return &WsClients{
		clientmaps: make(map[string]Client),
		mutex:      &sync.RWMutex{},
	}
}

func (w *WsClients) Set(c Client) error {
	if c.Id == "" {
		return errors.New("id cannot be empty")
	}
	if _, ok := w.Get(c.Id); ok {
		return errors.New("duplicate id")
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.clientmaps[c.Id] = c

	return nil
}

func (w *WsClients) Range(fn func(Client)) {
	w.mutex.RLock()
	for _, c := range w.clientmaps {
		go fn(c)
	}
	w.mutex.RUnlock()

}

func (w *WsClients) Get(id string) (Client, bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	client, ok := w.clientmaps[id]
	return client, ok
}

func (w *WsClients) Del(id string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	delete(w.clientmaps, id)
}
