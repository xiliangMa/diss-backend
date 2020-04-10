package ws

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	DissClient map[string]*Client

	// Register requests from the clients.
	register chan *Client
}

func NewHub() *Hub {
	return &Hub{
		DissClient: make(map[string]*Client),
		register:   make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.DissClient[client.clientIp] = client
		}
	}
}
