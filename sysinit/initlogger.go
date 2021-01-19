package sysinit

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
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
	FileName, _ := web.AppConfig.String("logs::FileName")
	MaxLines, _ := web.AppConfig.Int64("logs::MaxLines")
	MaxSize, _ := web.AppConfig.Int64("logs::MaxSize")
	Daily, _ := web.AppConfig.Bool("logs::Daily")
	MaxDays, _ := web.AppConfig.Int64("logs::MaxDays")
	Level, _ := web.AppConfig.Int64("logs::Level")
	Color, _ := web.AppConfig.Bool("logs::Color")

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
	logs.SetLogger(logs.AdapterMultiFile, string(jsonbyte))
}
