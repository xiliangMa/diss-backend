package ws

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/ws"
	"net/http"
)

type WSMetricController struct {
	web.Controller
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this *WSMetricController) Metrics() {
	wsm := &models.WSManager{
		Request:  this.Ctx.Request,
		Response: this.Ctx.ResponseWriter,
	}
	// 创建全局ws控制对象
	wsm.NewWSManager()
	err, wsconn := wsm.Err, wsm.Conn
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
		go wsmh.Save()
	}
	this.Data["json"] = nil
	this.ServeJSON(false)
}
