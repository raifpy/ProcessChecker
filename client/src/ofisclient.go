package ofisclient

import (
	"client/src/dialog"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
)

type OfisClient struct {
	Config   Config
	Conn     *websocket.Conn
	Hostname string
}

func NewOfisClient(c Config) (o *OfisClient, err error) {
	if c.Context == nil {
		c.Context = context.Background()
	}

	host, err := Hostname()
	if err != nil {
		return nil, err
	}

	fmt.Printf("host: %v\n", host)
	conn, response, err := websocket.DefaultDialer.DialContext(c.Context, c.DialAddr, http.Header{
		"id": {host},
	})
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err
	}

	return &OfisClient{
		Hostname: host,
		Config:   c,
		Conn:     conn,
	}, nil
}

func (o *OfisClient) Reconnect() error {
	conn, resp, err := websocket.DefaultDialer.Dial(o.Config.DialAddr, http.Header{
		"id": {o.Hostname},
	})
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return err
	}

	o.Conn = conn

	return nil
}

func (o *OfisClient) Run() error {
	log.Println("running")
	for {

		if err := o.Listen(); err != nil {

			if err := o.Reconnect(); err != nil {
				dialog.PopUp("hata: " + err.Error())
				time.Sleep(time.Second * 5)
				continue
			} else {
				dialog.PopUp("Yeniden bağlantı kuruldu!")
			}

			if err := o.Listen(); err != nil {
				dialog.PopUp("hata (listen): " + err.Error())
				time.Sleep(time.Second * 5)
				continue
			} else {
				dialog.PopUp("Yeniden bağlantı kuruldu!")
			}
		}
	}
}
