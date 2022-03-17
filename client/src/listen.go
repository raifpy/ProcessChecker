package ofisclient

func (o *OfisClient) Listen() error {
	for {
		var rq WebsocketResponse
		if err := o.Conn.ReadJSON(&rq); err != nil {
			return err
		}

		go o.responseParser(rq)

	}

}
