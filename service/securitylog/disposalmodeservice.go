package securitylog

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
)

type DisposalModeService struct {
}

func (this *DisposalModeService) DisposalMode(dm *models.DisposalMode) models.Result {
	natsManager := models.Nats
	result := models.Result{}
	if natsManager != nil && natsManager.Conn != nil {

		if dm.Action == models.WarnWhiteListConfigKey {
			for _, w := range dm.WarningWhiteList {
				whiteListConfig := w
				whiteListConfig.Add()

				warningInfo := new(models.WarningInfo)
				warningInfo.Mode = dm.Action
				warningInfo.Id = whiteListConfig.WarningInfoId
				warningInfo.Analysis = whiteListConfig.Desc
				warningInfo.Status = models.WarnInfoStatus
				warningInfo.Update()
			}
		} else {
			for _, wi := range dm.WarningInfo {
				warningInfo := wi
				if dm.Action != "" {
					warningInfo.Action = strings.Title(dm.Action + models.Container)
					warningInfo.Mode = dm.Action

					result.Code = http.StatusOK
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
					containerRespCenter.Status = "执行中"
					containerRespCenter.Add()

				}
				warningInfo.Status = models.WarnInfoStatus
				warningInfo.Update()
			}
		}
	}
	result.Code = http.StatusOK
	return result
}
