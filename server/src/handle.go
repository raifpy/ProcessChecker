package ofis

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/src/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/cast"
)

func (o *Ofis) Router() error {
	o.Engine.Use(logger.New())
	o.Engine.Get("/websocket", func(c *fiber.Ctx) error {

		id := c.Get("id")
		if id == "" {
			return fiber.ErrUnauthorized
		}

		if _, ok := o.Ws.Clients.Get(id); ok {
			return errors.New("multiple connection not allowed")
		}
		if !o.Ws.IsWebSocketUpgrade(c) {
			return fiber.ErrBadGateway
		}
		c.Locals("id", id)
		return c.Next()

	}).Get("/websocket", o.Ws.Middleware())

	o.Engine.Get("/sorgu", func(c *fiber.Ctx) error {
		processname := c.Query("process")
		sorgulist := []websocket.Sorgu{}
		for _, l := range o.Ws.WriteAll(websocket.ClientRequest{
			Message: websocket.RequestMessageProcessCheck{
				ProcessName: processname,
			}}) {

			oneymis := l.Message.(websocket.ClientParsedResponse)
			fmt.Printf("oneymis: %v\n", oneymis)

			fmt.Printf("l.Message: %v\n", l.Message)

			v, _ := json.Marshal(l.Message)
			fmt.Printf("v: %v\n", string(v))

			fmt.Printf("l.Error: %v\n", l.Error)
			sorgulist = append(sorgulist, websocket.Sorgu{
				Hostname: l.Id,
				Process:  cast.ToBool(oneymis.Message),
				Error:    cast.ToString(l.Error),
			})
		}

		return c.JSON(sorgulist)

	})

	return nil
}

func (o *Ofis) Listen() error {
	return o.Engine.Listen(o.Config.ListenAddr)
}
