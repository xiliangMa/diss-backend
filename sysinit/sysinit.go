package sysinit

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/k8s"
	"github.com/xiliangMa/diss-backend/service/nats"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/service/task"
	"github.com/xiliangMa/diss-backend/sysinit/db"
)

func init() {
	// initDB
	defaultDB := db.DefaultDB{}
	defaultDB.InitDB()
	//
	//securityLogDb := db.SecurityLogDb{}
	//securityLogDb.InitDB()
	//
	//dissApiDB := db.DissApiDB{}
	//dissApiDB.InitDB()

	//init logger
	InitLogger()

	// init global logConfig
	InitGlobalLogConfig()

	// init global syslog handler
	task.InitGlobalSyslogHander()

	// init trial license
	InitTrialLicense()

	// ================= 初始化全局变量
	// 1. NatsManager
	models.Nats = models.NewNatsManager()
	nats.RunClientSub_Image_Safe()
	nats.RunClientSub_IDL()

	// 2. WSManager
	models.WSHub = models.NewHub()
	go models.WSHub.Run()

	// 3. GRManager
	models.GRM = models.NewGoRoutineManager()

	// 4. KCM
	models.KCM = models.NewKubernetesClientManager()
	models.InitClientHub()
	models.InitDymaicClientHub()

	// watch all cluster
	k8sWatchService := k8s.K8sWatchService{}
	k8sWatchService.WatchAll()

	// init task
	InitTask()

	// 5. MailServerManager
	models.MSM = models.NewMailServerManager()
	mailService := system.MailService{}
	go mailService.StartMailService()

	// 6. LDAPManager
	models.LM = models.NewLDAPManager()

	//7. init AnchoreEngineManager
	models.AEM = models.NewAnchoreEngineManager()
}
