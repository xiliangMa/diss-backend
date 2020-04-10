package ws

import (
	"github.com/astaxie/beego/context"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

type WSManager struct {
	Conn *websocket.Conn
	Err  error
}

func (this *WSManager) NewWSManager(response *context.Response, request *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	this.Conn, this.Err = upgrader.Upgrade(response, request, nil)

	if this.Err == nil {
		ip := strings.Split(this.Conn.RemoteAddr().String(), ":")
		client := &Client{hub: WSHub, conn: this.Conn, send: make(chan []byte, 256), clientIp: ip[0]}
		client.hub.register <- client
		WS = this
	}
}

func (this *WSManager) GetWSManager() *WSManager {
	return WS
}
