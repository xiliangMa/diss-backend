package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WSMetricController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this *WSMetricController) Metrics() {
	wsconn, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		logs.Info("upgrade:", err)
		return
	}
	defer wsconn.Close()
	for {
		mt, message, err := wsconn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		logs.Info("recv: %s", message)

		err = wsconn.WriteMessage(mt, message)
		if err != nil {
			logs.Info("write:", err)
			break
		}
	}
	this.Data["json"] = nil
	this.ServeJSON(false)
}
