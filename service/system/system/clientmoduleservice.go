package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

type ClientModuleService struct {
	ConfigForUpdate     *models.UpdateAssets
	ClientModuleControl *models.ClientModuleControl
	HostIds             string
}

func (this *ClientModuleService) SetModule() *models.NatsData {
	cmsdata := this.ClientModuleControl
	if cmsdata != nil && cmsdata.ModuleName != "" && this.HostIds != "" {
		subject := this.HostIds
		result := models.NatsData{Type: models.Type_Control, Tag: models.Resource_ClientModuleControl, Data: cmsdata, RCType: models.Resource_Control_Type_Post}
		data, _ := json.MarshalIndent(result, "", "  ")
		err := models.Nats.Conn.Publish(subject, data)
		if err == nil {
			logs.Info("Deliver ClientModuleControl to Nats Success, Subject: %s Id: %s, data: %v", subject, cmsdata.ModuleName, result)
		}
		return &result
	}

	return nil
}

func (this *ClientModuleService) GetUpdateAssetsConfig() []byte {
	cmsdata := this.ConfigForUpdate

	if cmsdata != nil {
		result := models.NatsData{Type: models.Type_Control, Tag: models.Resource_AssetData, Data: cmsdata, RCType: models.Resource_Control_Type_Post}
		data, _ := json.MarshalIndent(result, "", "  ")

		return data
	}

	return nil
}
