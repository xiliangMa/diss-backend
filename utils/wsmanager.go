package utils

import (
	"github.com/astaxie/beego/context"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	WS *WSManager
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
	WS = this
}

func (this *WSManager) GetWSManager() *WSManager {
	return WS
}
