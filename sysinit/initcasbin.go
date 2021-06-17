package sysinit

import (
	"github.com/beego/beego/adapter/logs"
	"github.com/casbin/casbin"
	xormadapter "github.com/casbin/xorm-adapter"
	"github.com/xiliangMa/diss-backend/sysinit/db"
	"github.com/xiliangMa/diss-backend/utils"
)

func InitCasbinAdaptor() *casbin.Enforcer {
	DSAlias := utils.DS_Default
	DS := db.GetConn(DSAlias)
	adaptor, err := xormadapter.NewAdapter("postgres", DS)
	if err != nil {
		logs.Warn("create casbin adaptor error", err)
	}
	enforcer, err2 := casbin.NewEnforcer("conf/rbac_model.conf", adaptor)
	if err2 != nil {
		logs.Warn("create casbin enforce error", err2)
	}

	return enforcer
}
