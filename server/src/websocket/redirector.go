package websocket

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Redirector struct {
	Timeout    time.Duration
	redirector map[string]chan ClientParsedResponse
	mutex      *sync.RWMutex
}

func (r *Redirector) Get(id string) (ch chan ClientParsedResponse, ok bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ch, ok = r.redirector[id]
	return ch, ok
}

func (r *Redirector) Set(id string, ch chan ClientParsedResponse) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.redirector[id] = ch
}

func (r *Redirector) SetAsTimeout(id string, ch chan ClientParsedResponse, duration time.Duration) {
	r.Set(id, ch)
	time.AfterFunc(duration, func() {
		if _, ok := r.Get(id); ok {
			log.Println("Timeout exceeded: ", id)
			ch <- ClientParsedResponse{
				Error: errors.New("timeout exceeded"),
			}
			//r.Del(id)
		}
	})
}

func (r *Redirector) Del(id string) {
	if _, ok := r.Get(id); ok {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		delete(r.redirector, id)
	}
}

func (r *Redirector) Close(id string) {
	defer func() {
		_ = recover()
	}()

	ch, ok := r.Get(id)
	if !ok {
		return
	}

	r.Del(id)
	close(ch)
}

func NewRedirector() *Redirector {
	return &Redirector{
		Timeout:    time.Second * 10,
		redirector: make(map[string]chan ClientParsedResponse),
		mutex:      &sync.RWMutex{},
	}
}
