// @APIVersion 1.0.0
// @Title DISS API
// @Description DISS API
// @Schemes [http]

package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/xiliangMa/diss-backend/controllers"
	"github.com/xiliangMa/diss-backend/utils"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/hosts",
			beego.NSInclude(
				&controllers.HostController{},
			),
		),
	)

	var isLogin = func(ctx *context.Context) {
		if ctx.Request.Method != "OPTIONS" {
			_, code := utils.CheckToken(ctx.Input.Header("token"))
			if code != 200 {
				ctx.Redirect(401, "/swagger")
			}
		}

	}
	beego.InsertFilter("/v1/hosts/*", beego.BeforeRouter, isLogin)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.AddNamespace(ns)
}
