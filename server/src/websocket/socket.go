package websocket

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (s *Ws) IsWebSocketUpgrade(c *fiber.Ctx) bool {
	return websocket.IsWebSocketUpgrade(c)
}

func (s *Ws) Middleware() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		id, ok := c.Locals("id").(string)
		if !ok || id == "" {
			log.Println("Conn have incorrect id!")
			c.Close()
		}

		log.Printf("New Conn connected! %s\n", id)
		if err := s.Clients.Set(Client{
			Id:   id,
			Conn: c,
		}); err != nil {
			log.Println("Id set error: ", err)
			c.WriteJSON(ClientRequest{
				MessageID: "system:error",
				Message: MessageError{
					Error: err.Error(),
				},
			})
			c.Close()
			return
		}
		defer s.Clients.Del(id)

		for {
			var template ClientResponse
			if err := c.ReadJSON(&template); err != nil {
				log.Printf("Conn %s read error: %v\n", id, err)
				break
			}
			if template.ReplyID == "" {
				log.Printf("Conn %s readed value unsigned!\n", id)
				continue
			}
			log.Printf("Conn %s recevied %s", id, template.ReplyID)
			ch, ok := s.Redirector.Get(template.ReplyID)
			if !ok {
				log.Printf("Conn %s readed value id (%s) unsigned!\n", id, template.ReplyID)
				continue
			}

			go func() {
				ch <- ClientParsedResponse{
					Id:      id,
					Time:    time.Now(),
					Message: template.Message,
				}
			}()

		}

		log.Println("Conn disconnected! ", id)

	})
}
