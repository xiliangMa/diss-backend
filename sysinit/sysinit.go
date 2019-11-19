package sysinit

import "github.com/xiliangMa/diss-backend/utils"

func init() {
	InitDB()
	InitLogger()
	utils.InitEsClient()
}
