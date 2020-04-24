package ws

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net/http"
)

type WSManager struct {
	Conn     *websocket.Conn
	Err      error
	Response *context.Response
	Request  *http.Request
}

func (this *WSManager) NewWSManager() *WSManager {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	this.Conn, this.Err = upgrader.Upgrade(this.Response, this.Request, nil)
	if this.Err != nil {
		logs.Error("Create WSManager fail, err: %s", this.Err)
	}
	return this
}
