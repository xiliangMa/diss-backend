package sysinit

import (
	"github.com/xiliangMa/diss-backend/models"
)

func InitCasbin() {
	models.GlobalCasbin = models.NewCasbinManager()
	models.GlobalCasbin.Enforcer.LoadPolicy()
}
