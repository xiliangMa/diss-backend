package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type HostService struct {
	Host models.HostConfig
}

// todo add result
func (this *HostService) Delete() models.Result {
	result := models.Result{Code: http.StatusOK}
	hostId := this.Host.Id

	// delete ImageConfig
	ic := models.ImageConfig{}
	ic.HostId = hostId
	result = ic.Delete()

	// delete ImageInfo
	ii := models.ImageInfo{}
	ii.HostId = hostId
	result = ii.Delete()

	// delete ContainerConfig
	cc := models.ContainerConfig{}
	cc.HostId = hostId
	result = cc.Delete()

	// delete ContainerInfo
	ci := models.ContainerInfo{}
	ci.HostId = hostId
	result = ci.Delete()

	// delete HostInfo
	hi := models.HostInfo{}
	hi.Id = hostId
	result = hi.Delete()

	// delete HostConfig
	result = this.Host.Delete()
	return result
}

func (this *HostService) SendAuthorizationDetail() models.Result {
	result := models.Result{}
	hc := this.Host.Get()

	if hc != nil {
		result.Code = http.StatusOK
		subject := hc.Id
		r := models.NatsData{Code: result.Code, Type: models.Type_Config, Msg: result.Message, Tag: models.Resource_Authorization, RCType: models.Resource_Control_Type_Put, Data: hc}
		data, _ := json.MarshalIndent(r, "", "  ")
		err := models.Nats.Conn.Publish(subject, data)
		if err != nil {
			logs.Error("Nats ############################ received data from agent fail ############################", err)
		}
	} else {
		logs.Error("Host Authorization fail")
	}
	return result
}
