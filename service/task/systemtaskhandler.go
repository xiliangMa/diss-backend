package task

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"time"
)

type SystemCheckHandler struct {
}

func (this *SystemCheckHandler) SystemCheck() {
	// host 相关状态检查
	data := new(models.HostConfig).List(0, 0)
	items := data.Data.(map[string]interface{})["items"]
	for _, host := range items.([]*models.HostConfig) {
		now := time.Now().UTC()
		sub := int(now.Sub(host.UpdateTime).Minutes())
		agentHeartBeatTime, _ := beego.AppConfig.Int("system::AgentHeartBeatTime")
		if sub > agentHeartBeatTime {
			host.Diss = models.Diss_NotInstalled
			host.DissStatus = models.Diss_Status_Unsafe
			host.Status = models.Host_Status_Abnormal
			host.Update()
			logs.Warn("Heartbeat abnormal, HostId: %s, Duration: %v Minutes", host.Id, sub)
		} else {
			host.Diss = models.Diss_Installed
			host.DissStatus = models.Diss_status_Safe
			host.Status = models.Host_Status_Normal
			host.Update()
		}
	}
}
