package websocket

type Ws struct {
	Redirector *Redirector
	Clients    *WsClients
}

func NewWs() (ws *Ws, err error) {
	ws = &Ws{
		Redirector: NewRedirector(),
		Clients:    NewWsClients(),
	}

	return
}
