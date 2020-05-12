package sysinit

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/sysinit/db"
)

func init() {
	// initDB
	defaultDB := db.DefaultDB{}
	defaultDB.InitDB()

	//securityLogDb := db.SecurityLogDb{}
	//securityLogDb.InitDB()
	//
	//dissApiDB := db.DissApiDB{}
	//dissApiDB.InitDB()

	//init logger
	InitLogger()

	// init task
	InitTask()

	// ================= 初始化全局变量
	// 1. NatsManager
	models.Nats = models.NewNatsManager()

	// 2. WSManager
	models.WSHub = models.NewHub()
	go models.WSHub.Run()
}
