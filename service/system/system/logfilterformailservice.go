package system

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"strings"
)

type LogToMailFilterService struct {
	InfoType    string
	InfoSubType string
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
	var HighTags = models.WarningLevel_High
	var MediumTags = models.WarningLevel_Medium
	var LowTags = models.WarningLevel_Low
	var level = models.BML_Level_ALL

	if this.InfoSubType == models.WarningInfo_Image {
		// 镜像安全类型转换
		this.LogType = models.LogType_ImageSecLog
	} else if strings.HasPrefix(this.InfoSubType, "ALERT_TYPE") {
		// 入侵检测类型转换
		this.LogType = models.LogType_IntrudeDetectLog
	} else if this.InfoSubType == models.TMP_Type_BM_Docker || this.InfoType == models.TMP_Type_BM_K8S {
		// 基线类型不通过告警
		return level
	} else {
		this.LogType = this.InfoType
		HighTags = models.BML_Level_High
		MediumTags = models.BML_Level_Medium
		LowTags = models.BML_Level_Low
	}

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

func (this *LogToMailFilterService) SendToChannel(dataModel interface{}) {
	dataModelString, err := utils.ToIndentJSON(&dataModel)

	if err != nil {
		logs.Error("Encode json for mail fail, error: ", err)
	}
	dataModelMap := map[string]string{}
	dataModelMap[models.MailField_Subject] = models.LogToEmail_Prefix
	dataModelMap[models.MailField_LogType] = this.LogType
	dataModelMap[models.MailField_InfoSubType] = this.InfoSubType
	dataModelMap[models.MailField_Body] = "<pre>" + dataModelString + "</pre>"
	models.MSM.LogChannel <- &dataModelMap
}
