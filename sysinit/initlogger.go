package sysinit

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

type adapterMultiFileConfig struct {
	FileName string
	MaxLines int64
	MaxSize  int64
	Daily    bool
	MaxDays  int64
	Level    int64
	Color    bool
}

/***
 *
 *	LevelEmergency = iota
 *	LevelAlert
 *	LevelCritical
 *	LevelError
 *	LevelWarning
 *	LevelNotice
 *	LevelInformational
 *	LevelDebug
 *
 */
func InitLogger() {
	FileName := beego.AppConfig.String("logs::FileName")
	MaxLines, _ := beego.AppConfig.Int64("logs::MaxLines")
	MaxSize, _ := beego.AppConfig.Int64("logs::MaxSize")
	Daily, _ := beego.AppConfig.Bool("logs::Daily")
	MaxDays, _ := beego.AppConfig.Int64("logs::MaxDays")
	Level, _ := beego.AppConfig.Int64("logs::Level")
	Color, _ := beego.AppConfig.Bool("logs::Color")

	logConfig := adapterMultiFileConfig{
		FileName: FileName,
		MaxLines: MaxLines,
		MaxSize:  MaxSize,
		Daily:    Daily,
		MaxDays:  MaxDays,
		Level:    Level,
		Color:    Color,
	}

	jsonbyte, _ := json.MarshalIndent(logConfig, "", "")
	logs.Info("log config info %s", string(jsonbyte))
	beego.SetLogger(logs.AdapterMultiFile, string(jsonbyte))
}

func InitGlobalLogConfig() {
	var logConfig models.LogConfig
	logConfig.ConfigName = models.Log_Config_SysLog_Export
	syslogConfig := logConfig.InnerGet()
	if syslogConfig != nil {
		models.GlobalLogConfig[models.Log_Config_SysLog_Export] = syslogConfig[0]
	} else {
		models.GlobalLogConfig[models.Log_Config_SysLog_Export] = &logConfig
	}

}
