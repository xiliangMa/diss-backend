package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"strings"
)

type LogToMailFilterService struct {
	WarningType string
	LogType     string
	CurLevelTag string
}

func (this *LogToMailFilterService) GetLogFilterStatus() bool {

	logToMailConfig := models.MSM.LogToMailConfig

	level := this.TransLogTypeAndTag()

	filterStatus := false
	if logToMailConfig == nil {
		return filterStatus
	}

	curLogConfig := (*logToMailConfig)[this.LogType]
	if curLogConfig == nil {
		return filterStatus
	}

	if curLogConfig["Enable"] {
		// 先检查开关，再对比当前的LevelTag和各级别，确定是否放行
		if curLogConfig["Level_High"] && level == models.BML_Level_High {
			filterStatus = true
		} else if curLogConfig["Level_Medium"] && level == models.BML_Level_Medium {
			filterStatus = true
		} else if curLogConfig["Level_Low"] && level == models.BML_Level_Low {
			filterStatus = true
		}
	}

	return filterStatus
}

func (this *LogToMailFilterService) TransLogTypeAndTag() string {
	// 处理告警类型对应到日志类型的转换，以及对应设置等级标记
	// 默认设置的级别Tags 为 High Medium Low
	var HighTags = models.BML_Level_High
	var MediumTags = models.BML_Level_Medium
	var LowTags = models.BML_Level_Low

	if this.WarningType == models.WarningInfo_Image {
		// 镜像安全类型转换
		this.LogType = models.LogType_ImageSecLog
	} else if strings.HasPrefix(this.WarningType, "ALERT_TYPE") {
		// 入侵检测类型转换
		this.LogType = models.LogType_IntrudeDetectLog
		HighTags = models.WarningLevel_High
		MediumTags = models.WarningLevel_Medium
		LowTags = models.WarningLevel_Low
	} else if this.WarningType == models.TMP_Type_BM_Docker || this.WarningType == models.TMP_Type_BM_K8S {
		// 基线类型
		this.LogType = models.LogType_BenchMarkLog
		HighTags = models.WarningLevel_High
		MediumTags = models.WarningLevel_Medium
		LowTags = models.WarningLevel_Low
	} else {
		// 其他类型，如审计
		this.LogType = this.WarningType
	}

	var level = models.BML_Level_ALL
	if strings.Contains(HighTags, this.CurLevelTag) {
		//等级序号 高
		level = models.BML_Level_High
	} else if strings.Contains(MediumTags, this.CurLevelTag) {
		//等级序号 中
		level = models.BML_Level_Medium
	} else if strings.Contains(LowTags, this.CurLevelTag) {
		//等级序号 低
		level = models.BML_Level_Low
	}
	return level
}

func (this *LogToMailFilterService) SendToChannel(warningInfo models.WarningInfo) {
	warningInfoJson, err := json.Marshal(warningInfo)
	if err != nil {
		logs.Error("Encode warningInfo for mail json fail, error: ", err)
	}
	warningString := string(warningInfoJson)
	warningMap := map[string]string{}
	warningMap[models.MailField_Subject] = models.LogToEmail_Prefix
	warningMap[models.MailField_LogType] = this.LogType
	warningMap[models.MailField_Body] = warningString
	models.MSM.LogChannel <- &warningMap
}
