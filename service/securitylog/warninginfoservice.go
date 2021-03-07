package securitylog

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type WarningInfoService struct {
}

func (this *WarningInfoService) DisposalMode(warningInfo *models.WarningInfo) models.Result {
	natsManager := models.Nats
	result := models.Result{}
	if natsManager != nil && natsManager.Conn != nil {

		if warningInfo.Action != "" {
			warningInfo.Action = strings.Title(warningInfo.Action + models.Container)

			result.Code = http.StatusOK
			result.Message = "DisposalMode Success"

			subject := warningInfo.HostId
			r := models.NatsData{Code: result.Code, Type: models.Type_Control, Msg: result.Message, Tag: models.Resource_ContainerControl, RCType: models.Resource_Control_Type_Put, Data: warningInfo}
			data, _ := json.MarshalIndent(r, "", "  ")
			logs.Info("Publish success, subject: %s data : %s", subject, data)
			err := models.Nats.Conn.Publish(subject, data)
			if err != nil {
				logs.Error("Nats ############################ received data from agent fail ############################", err)
			}
		}

		warningInfo.Status = models.WarnInfoStatus
		res := warningInfo.Update()

		if res.Code == http.StatusOK {
			containerRespCenter := new(models.RespCenter)
			_ = utils.ConvertType(warningInfo, containerRespCenter)
			containerRespCenter.Id = ""
			containerRespCenter.WarningInfoId = warningInfo.Id
			containerRespCenter.Status = ""
			containerRespCenter.Add()
		}
	}

	return result
}
