package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

type SysConfigService struct {
	SysConfig *models.SysConfig
	ConfigKey string
	ParamName string
}

func (this *SysConfigService) RefreshConnections() {
	if this.SysConfig.Key == models.EmailServerConfig {
		models.MSM.NewMailDialer()
	}
	if this.SysConfig.Key == models.Log_Config_To_Mail {
		models.MSM.NewLogToMail()
	}
	if this.SysConfig.Key == models.LDAPClientConfig {
		models.LM.CreateLDAPConnection()
	}
}

func (this *SysConfigService) GetConfigDataBool() (status bool, err error) {
	sysConfig := models.SysConfig{}
	sysConfig.Key = this.ConfigKey
	sysConfigData := sysConfig.Get()
	if sysConfigData != nil && sysConfigData.Value != "" {
		sysConfigJson := map[string]bool{}
		err := json.Unmarshal([]byte(sysConfigData.Value), &sysConfigJson)
		if err != nil {
			logs.Error("Parse %s Config error %s.", models.Resource_WarningInfo, err)
			return false, err
		}
		paramName := this.ParamName
		return sysConfigJson[paramName], err
	}
	return false, err
}
