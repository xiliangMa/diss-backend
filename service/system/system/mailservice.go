package system

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"gopkg.in/gomail.v2"
	"net/http"
	"strings"
)

type MailService struct {
	Subject          string
	Body             string
	WarningInfoName  string
	WarningInfoError string
}

func (this *MailService) SendMail() error {

	var err error = nil
	dialer := models.MSM.Dialer

	if dialer != nil {
		// 邮箱服务可用，执行邮件发送
		mailConfig := models.MSM.MailServerConfig
		m := gomail.NewMessage()

		mailFrom := fmt.Sprintf("%s<%s>", mailConfig.MailFrom, mailConfig.UserName)
		m.SetHeader(models.MailField_From, mailFrom)
		mailTo := fmt.Sprintf("%s<%s>", mailConfig.MailTo, mailConfig.SendToEmail)
		m.SetHeader(models.MailField_To, mailTo)
		m.SetHeader(models.MailField_Subject, this.Subject)
		m.SetBody("text/html", this.Body)

		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		this.Body = strings.Replace(this.Body, "Action", "Sensi1", 1)
		this.Body = strings.Replace(this.Body, "Actor", "Sensi2", 1)
		err = dialer.DialAndSend(m)
	}

	if err != nil {
		logs.Error(err)
		return err
	}

	return err
}

func (this *MailService) StartMailService() {
	for {
		select {
		case message, ok := <-models.MSM.LogChannel:
			if ok {
				dialer := models.MSM.Dialer

				if dialer == nil {
					//记录邮箱不可用错误告警
					this.WarningInfoName = models.MailServer_Not_Available_Msg
					this.WarningInfoError = fmt.Sprint("MailServer Not Available , code %s", utils.MailServerNotAvaiableErr)
					warningInfo := this.FillMailWarningData()

					if result := warningInfo.Add(); result.Code != http.StatusOK {
						logs.Warn("MailServer not available, Error:", this.WarningInfoError)
					}
					break
				}
				this.Subject = (*message)[models.MailField_Subject] + " - " + (*message)[models.MailField_LogType] + " - " + (*message)[models.MailField_InfoSubType]
				this.Body = (*message)[models.MailField_Body]
				err := this.SendMail()

				if err != nil {
					// 记录邮件发送失败告警
					models.MSM.Status = models.MailServerStatus_Abnormal //标记发送失败
					this.WarningInfoName = models.Mail_CanNotSend_Msg
					this.WarningInfoError = err.Error()
					warningInfo := this.FillMailWarningData()

					if result := warningInfo.Add(); result.Code != http.StatusOK {
						logs.Warn("Send Email fail, Error:", err)
					}
				} else {
					models.MSM.Status = models.MailServerStatus_Normal
				}

			}
		}
	}
}

func (this *MailService) FillMailWarningData() models.WarningInfo {
	warningInfo := models.WarningInfo{}
	warningInfo.Type = models.WarningInfo_MailError
	warningInfo.Level = models.WarningLevel_Medium
	warningInfo.Status = models.WarningStatus_Not_Dealed
	warningInfo.Name = this.WarningInfoName
	warningInfo.Account = models.Account_Admin

	basicdata := map[string]string{
		models.MailField_Subject: this.Subject,
		models.MailField_From:    fmt.Sprintf("%s<%s>", models.MSM.MailServerConfig.MailFrom, models.MSM.MailServerConfig.UserName),
		models.MailField_To:      fmt.Sprintf("%s<%s>", models.MSM.MailServerConfig.MailTo, models.MSM.MailServerConfig.SendToEmail),
	}
	basicJson, _ := json.Marshal(basicdata)
	detail := map[string]string{
		"error": this.WarningInfoError,
	}
	detailsJson, _ := json.Marshal(detail)
	conjunctionJson := fmt.Sprintf(`{"basic": %s,"details": %s}`, basicJson, detailsJson)
	warningInfo.Info = string(conjunctionJson)
	return warningInfo
}
