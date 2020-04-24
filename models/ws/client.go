package ws

import "github.com/gorilla/websocket"

type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	ClientIp string

	SystemId string
}
