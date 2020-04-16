package ws

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/service/ws"
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
	wsm := new(ws.WSManager)
	wsm.Request = this.Ctx.Request
	wsm.Response = this.Ctx.ResponseWriter
	// 创建全局ws控制对象
	wsm.NewWSManager()
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
		wsmh := &ws.WSMetricsService{message, wsconn}
		wsmh.Save()
	}
	this.Data["json"] = nil
	this.ServeJSON(false)
}
