package system

import "github.com/xiliangMa/diss-backend/models"

type SysConfigService struct {
	SysConfig *models.SysConfig
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
