package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/beego/beego/v2/core/logs"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter"
	"github.com/dgrijalva/jwt-go"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	_ "strings"
)

type CasbinManager struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinManager() *CasbinManager {
	DSAlias := utils.DS_Default
	DS := utils.GetConn(DSAlias)
	adaptor, err := xormadapter.NewAdapter("postgres", DS)
	if err != nil {
		logs.Error("Create casbin adapter failed, err: %s.", err)
		return nil
	}
	enforcer, err := casbin.NewEnforcer("conf/rbac_model.conf", adaptor)
	if err != nil {
		logs.Error("Create casbin enforcer failed, err: %s.", err)
		return nil
	}
	logs.Info("New casbin manager success.")
	return &CasbinManager{Enforcer: enforcer}
}

func (this *CasbinManager) CheckPermisson() beego.FilterFunc {
	return func(ctx *context.Context) {
		var secret []byte = []byte("secret")
		tokenStr := ctx.Request.Header.Get("token")
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		name := token.Claims.(jwt.MapClaims)["UserName"].(string)
		userRoleList, _ := GlobalCasbin.Enforcer.GetRolesForUser(name)
		isAllow := false
		for _, userRole := range userRoleList {
			// todo check path
			//path := strings.ToLower(ctx.Request.URL.Path)
			module := ctx.Request.Header.Get("module")
			status, _ := GlobalCasbin.Enforcer.Enforce(userRole, module, "-")
			if status {
				isAllow = true
				return
			}
		}
		if !isAllow {
			result := Result{Code: http.StatusUnauthorized, Message: "权限不足"}
			ctx.Output.JSON(result, false, false)
		}
	}
}
