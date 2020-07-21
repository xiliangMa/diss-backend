package sysinit

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/k8s"
	"github.com/xiliangMa/diss-backend/service/task"
	"github.com/xiliangMa/diss-backend/sysinit/db"
)

func init() {
	// initDB
	defaultDB := db.DefaultDB{}
	defaultDB.InitDB()

	securityLogDb := db.SecurityLogDb{}
	securityLogDb.InitDB()

	dissApiDB := db.DissApiDB{}
	dissApiDB.InitDB()

	//init logger
	InitLogger()

	// sync All Cluster
	//k8STaskHandler := task.K8sTaskHandler{}
	//k8STaskHandler.SyncAll()

	// init task
	InitTask()

	// init global logConfig
	InitGlobalLogConfig()

	// init global syslog handler
	task.InitGlobalSyslogHander()

	// init trial license
	InitTrialLicense()

	// ================= 初始化全局变量
	// 1. NatsManager
	models.Nats = models.NewNatsManager()

	// 2. WSManager
	models.WSHub = models.NewHub()
	go models.WSHub.Run()

	// 3. GRManager
	models.GRM = models.NewGoRoutineManager()

	// watch all cluster
	k8sWatchService := k8s.K8sWatchService{}
	k8sWatchService.WatchAll()
}
