// @APIVersion 1.0.0
// @Title DISS API
// @Description DISS API
// @Schemes http, https
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/xiliangMa/diss-backend/controllers"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/hosts",
			beego.NSInclude(
				&controllers.HostController{},
			),
		),
		beego.NSNamespace("/merticinfo",
			beego.NSInclude(
				&controllers.MetricController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
	)

	// add route for ws
	beego.Router("/metrics", &controllers.WSMetricController{}, "*:Metrics")

	var isLogin = func(ctx *context.Context) {
		if ctx.Request.Method != "OPTIONS" {
			_, code := utils.CheckToken(ctx.Input.Header("token"))
			if code != http.StatusOK {
				ctx.Redirect(http.StatusUnauthorized, "/swagger")
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
