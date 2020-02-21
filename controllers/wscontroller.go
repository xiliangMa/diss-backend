package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/service"
	"github.com/xiliangMa/diss-backend/utils"
)

type WSMetricController struct {
	beego.Controller
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
		mt, message, err := wsconn.ReadMessage()
		if err != nil {
			logs.Error("read:", err)
			break
		}
		wsmh := &service.WSMetricsService{message}
		wsmh.Save()

		logs.Info("recv: %s", message)

		err = wsconn.WriteMessage(mt, message)
		if err != nil {
			logs.Error("write:", err)
			break
		}
	}
	this.Data["json"] = nil
	this.ServeJSON(false)
}
