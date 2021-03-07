package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
	"strings"
)

type RespCenterService struct {
}

func (this *RespCenterService) ContainerOperation(crc *models.RespCenter) models.Result {

	natsManager := models.Nats
	result := models.Result{}
	if natsManager != nil && natsManager.Conn != nil {

		if crc.Action != "" {
			crc.Action = strings.Title(crc.Action + models.Container)

			result.Code = http.StatusOK
			result.Message = "Container Operation Success"

			subject := crc.HostId
			r := models.NatsData{Code: result.Code, Type: models.Type_Control, Msg: result.Message, Tag: models.Resource_ContainerControl, RCType: models.Resource_Control_Type_Put, Data: crc}
			data, _ := json.MarshalIndent(r, "", "  ")
			logs.Info("Publish suss, subject: %s data : %s", subject, data)
			err := models.Nats.Conn.Publish(subject, data)
			if err != nil {
				logs.Error("Nats ############################ received data from agent fail ############################", err)
			}
		}
		return crc.Update()

	}
	return result

}
