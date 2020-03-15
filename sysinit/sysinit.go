package sysinit

import "github.com/xiliangMa/diss-backend/sysinit/db"

func init() {
	db.InitDB()
	db.InitSecurityLogDB()
	db.InitDissApiDB()
	InitLogger()
	InitTask()
}
