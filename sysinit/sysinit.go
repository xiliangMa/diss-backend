package sysinit

import (
	"github.com/xiliangMa/diss-backend/models/global"
	modelws "github.com/xiliangMa/diss-backend/models/ws"
	"github.com/xiliangMa/diss-backend/service/nats"
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

	// init task
	InitTask()

	// ================= 初始化全局变量
	// 1. NatsManager
	global.NatsManager = nats.NewNatsManager()

	// 2. WSManager
	global.WSHub = modelws.NewHub()
	go global.WSHub.Run()
}
