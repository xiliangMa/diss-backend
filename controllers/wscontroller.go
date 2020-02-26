package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/service"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
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
	wsm := new(utils.WSManager)
	// 创建全局ws控制对象
	wsm.NewWSManager(this.Ctx.ResponseWriter, this.Ctx.Request)
	err, wsconn := wsm.GetWSManager().Err, wsm.GetWSManager().Conn
	if err != nil {
		logs.Info("upgrade:", err)
		return
	}
	defer wsconn.Close()
	for {
		_, message, err := wsconn.ReadMessage()
		if err != nil {
			logs.Info("############################ Sync agent data fail ############################, err: ", err)
			break
		}

		wsmh := &service.WSMetricsService{message}
		wsmh.Save()

		//err = wsconn.WriteMessage(mt, message)
		respmsg := "received ok: " + strconv.Itoa(len(message)) + " bytes "

		err = wsconn.WriteMessage(websocket.TextMessage, []byte(respmsg))
		if err != nil {
			logs.Info("############################ Received data from agent fail ############################", err)
			break
		}
	}
	this.Data["json"] = nil
	this.ServeJSON(false)
}
