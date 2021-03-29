package sysinit

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/service/task"
)

func InitSysLogConfig(resetime bool) {

	system.GlobalSyslog.UpdateSysLogConfCache()

	if resetime {
		InitTimeEdgePoint()
	}

}

func InitTimeEdgePoint() {

	for exType, exGroup := range *task.AllExGroups() {

		TEPoint := new(models.TimeEdgePoint)
		TEPoint.EdgePointCode = exType
		TEPinDB := TEPoint.Get()
		if len(TEPinDB) == 0 {
			//初始起始时间设置为2018-1-1
			uid, _ := uuid.NewV4()
			TEPoint.TimePointA = models.DefaultStartTimeStamp
			TEPoint.EdgePointName = exGroup.ExportName
			TEPoint.TEPointId = uid.String()
			TEPoint.Direction = "lookback"
			TEPoint.ScopeSymbol = "----|"
			TEPoint.Add()
		}
	}
}
