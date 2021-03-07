package securitylog

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type DisposalModeService struct {
}

func (this *DisposalModeService) DisposalMode(dm *models.DisposalMode) models.Result {
	natsManager := models.Nats
	result := models.Result{}
	if natsManager != nil && natsManager.Conn != nil {

		warningInfo := new(models.WarningInfo)

		if dm.WarningWhiteList != (models.WarningWhiteList{}) {

			whiteListConfig := dm.WarningWhiteList

			if dm.Action == models.WarnWhiteListConfigKey {
				warningInfo.Mode = dm.Action
				whiteListConfig.Add()
			}

			warningInfo.Id = whiteListConfig.WarningInfoId
			warningInfo.Proposal = whiteListConfig.Desc
			warningInfo.Status = models.WarnInfoStatus
		}

		if dm.WarningInfo != (models.WarningInfo{}) {

			warningInfo = &dm.WarningInfo

			if dm.Action != "" {

				warningInfo.Action = strings.Title(dm.Action + models.Container)

				warningInfo.Mode = dm.Action
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

				containerRespCenter := new(models.RespCenter)
				_ = utils.ConvertType(warningInfo, containerRespCenter)
				containerRespCenter.Id = ""
				containerRespCenter.WarningInfoId = warningInfo.Id
				containerRespCenter.Status = ""
				containerRespCenter.Add()
			} else {
				warningInfo.Status = models.WarnInfoStatus
			}
		}

		return warningInfo.Update()
	}

	return result
}
