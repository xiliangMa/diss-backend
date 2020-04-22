package sysinit

import "github.com/xiliangMa/diss-backend/sysinit/db"

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
}
