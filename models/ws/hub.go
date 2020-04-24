package ws

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	DissClient map[string]*Client

	// Register requests from the clients.
	Register chan *Client
}

func NewHub() *Hub {
	return &Hub{
		DissClient: make(map[string]*Client),
		Register:   make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.DissClient[client.SystemId] = client
		}
	}
}
