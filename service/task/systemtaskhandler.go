package task

import (
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

type SystemCheckHandler struct {
}

func (this *SystemCheckHandler) SystemCheck() {
	// host 相关状态检查
	data := new(models.HostConfig).List(0, 0)
	if data.Code != http.StatusOK || data.Data == nil {
		return
	}
	items := data.Data.(map[string]interface{})["items"]
	for _, host := range items.([]*models.HostConfig) {
		now := time.Now()
		nano := time.Unix(host.HeartBeat/1e9, 0)
		sub := int(now.Sub(nano).Minutes())
		agentHeartBeatTime, _ := beego.AppConfig.Int("system::AgentHeartBeatTime")
		if host.IsEnableHeartBeat {
			if sub >= agentHeartBeatTime {
				host.Diss = models.Diss_NotInstalled
				host.Status = models.Host_Status_Abnormal
				host.DissStatus = models.Diss_Status_Unsafe
				host.OfflineTime = host.HeartBeat
				host.Update()
				logs.Warn("Heartbeat abnormal, HostId: %s, Duration: %v Minutes", host.Id, sub)
			}
		}
	}
}
