package ofis

import (
	"server/src/websocket"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	ListenAddr string
}

type Ofis struct {
	Config Config
	Ws     *websocket.Ws
	Engine *fiber.App
}

func NewOfis(c Config) (o *Ofis, err error) {
	o = &Ofis{
		Config: c,
		Engine: fiber.New(fiber.Config{
			EnablePrintRoutes:     true,
			DisableStartupMessage: true,
		}),
	}

	if o.Ws, err = websocket.NewWs(); err != nil {
		return
	}

	err = o.Router()

	return

}
