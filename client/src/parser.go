package ofisclient

import (
	"client/src/dialog"
	"fmt"
	"log"
	"strings"
)

func (o *OfisClient) responseParser(response WebsocketResponse) {
	fmt.Printf("messageid: %v\n", response.MessageID)
	if response.MessageID == "" {
		log.Println("empty message id!")
		return
	}
	if strings.Contains(response.MessageID, "system") {
		fmt.Println("-system message-")

		split := strings.Split(response.MessageID, ":")
		if len(split) == 1 {
			return
		}

		m := fmt.Sprint(response.Message)
		fmt.Printf("system message: %v\n", m)

		switch split[1] {
		case "error":
			{
				go dialog.Error(m)
			}
		case "info":
			{
				go dialog.Info(m)
			}
		}

		return
	}

	resparsed, ok := response.Message.(map[string]interface{})
	if !ok {
		log.Printf("!ok %T\n", response.Message)
	}
	fmt.Printf("resparsed: %v\n", resparsed)

	if errs, ok := resparsed["Error"]; ok {
		dialog.Error(fmt.Sprint(errs))
		return
	} else if ProcessName, ok := resparsed["ProcessName"]; ok {

		o.Conn.WriteJSON(WebsocketRequest{
			ReplyID: response.MessageID,
			Message: CheckProcess(fmt.Sprint(ProcessName)),
		})
	}

}
