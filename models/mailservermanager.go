package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/gomail.v2"
)

type MailServerManager struct {
	Dialer           *gomail.Dialer
	MailServerConfig *MailServerConfig
	LogToMailConfig  *map[string]map[string]bool
	LogChannel       chan *map[string]string
	Status           string
}

type MailServerConfig struct {
	EmailServer     string
	EmailServerPort int
	UserName        string
	Password        string
	SendToEmail     string
	MailFrom        string
	MailTo          string
	EnableSendLog   bool
}

func NewMailServerManager() *MailServerManager {
	MSManager := &MailServerManager{
		LogChannel: make(chan *map[string]string),
	}
	MSManager.NewMailDialer()
	MSManager.NewLogToMail()
	return MSManager
}

func (this *MailServerManager) NewMailDialer() {
	sysConfig := SysConfig{}
	sysConfig.Key = EmailServerConfig
	mailServerConfigData := sysConfig.Get()
	if mailServerConfigData != nil {
		mailConfigStr := mailServerConfigData.Value
		mailServerConfig := MailServerConfig{}
		err := json.Unmarshal([]byte(mailConfigStr), &mailServerConfig)
		if err != nil {
			logs.Error("Encode mailConfig json fail, error: ", err)
		}

		this.MailServerConfig = &mailServerConfig
	}
	if this.MailServerConfig.SendToEmail != "" && this.MailServerConfig.UserName != "" && this.MailServerConfig.Password != "" {
		logs.Info("MailConfig From: %s <%s> to: %s <%s>", this.MailServerConfig.MailFrom, this.MailServerConfig.UserName, this.MailServerConfig.MailTo, this.MailServerConfig.SendToEmail)
		this.Dialer = gomail.NewDialer(this.MailServerConfig.EmailServer, this.MailServerConfig.EmailServerPort, this.MailServerConfig.UserName, this.MailServerConfig.Password)
	}
}

func (this *MailServerManager) NewLogToMail() {
	sysConfig := SysConfig{}
	sysConfig.Key = Log_Config_To_Mail
	logMailConfig := sysConfig.Get()
	if logMailConfig == nil {
		return
	}
	logConfigMap := make(map[string]map[string]bool)

	logToMail := logMailConfig.Value
	err := json.Unmarshal([]byte(logToMail), &logConfigMap)
	if err != nil {
		logs.Error("Encode logMailConfig json fail, error: ", err)
	}
	this.LogToMailConfig = &logConfigMap
}
