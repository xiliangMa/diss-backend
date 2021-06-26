package plugins

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/xiliangMa/diss-backend/service/auth"
	"net/http"
)

type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

func (a *BasicAuthorizer) GetUserRole(input *context.Context) string {

	jwtService := auth.JwtService{}
	jwtService.TokenStr = input.Request.Header.Get("token")
	username := jwtService.GetUserFromToken()

	return username
}

func NewAuthorizer(e *casbin.Enforcer) beego.FilterFunc {
	return func(ctx *context.Context) {
		// 检验token
		jwtService := auth.JwtService{}
		if ctx.Request.Method != "OPTIONS" {
			_, code := jwtService.CheckToken(ctx.Request.Header.Get("token"))
			if code != http.StatusOK {
				ctx.Redirect(http.StatusUnauthorized, "/swagger")
				return
			}

			// 授权检查
			module := ctx.Request.Header.Get("module")
			a := &BasicAuthorizer{enforcer: e}
			username := jwtService.Token.Claims.(jwt.MapClaims)["UserName"].(string)
			// 后续完善可以通过path转换为模块进行验证
			// path := strings.ToLower(ctx.Request.URL.Path)

			userRoleList, _ := a.enforcer.GetRolesForUser(username)

			isAllow := false
			for _, userRole := range userRoleList {
				if status, err := a.enforcer.Enforce(userRole, module, "-"); status {
					if err != nil {
						logs.Info("Get enforce with role fail, error ", err)
					} else {
						isAllow = true
					}
				}
			}
			if !isAllow {
				ctx.Output.Status = 401
				msgNoPerm := map[string]string{"msg": "用户权限不足"}
				_ = ctx.Output.JSON(msgNoPerm, beego.BConfig.RunMode != "prod", false)
			}

		}
	}
}
