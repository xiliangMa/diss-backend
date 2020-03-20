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
	caccounts "github.com/xiliangMa/diss-backend/controllers/accounts"
	ca "github.com/xiliangMa/diss-backend/controllers/asset"
	cgroups "github.com/xiliangMa/diss-backend/controllers/groups"
	csl "github.com/xiliangMa/diss-backend/controllers/securitylog"
	cs "github.com/xiliangMa/diss-backend/controllers/securitypolicy"
	cstatistics "github.com/xiliangMa/diss-backend/controllers/statistics"
	css "github.com/xiliangMa/diss-backend/controllers/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1/statistics",
			beego.NSInclude(
				&cstatistics.StatisticsController{},
			),
		),
		beego.NSNamespace("/v1/hosts",
			beego.NSInclude(
				&controllers.HostController{},
			),
		),
		beego.NSNamespace("/v1/asset/images",
			beego.NSInclude(
				&ca.ImageController{},
			),
		),
		beego.NSNamespace("/v1/asset/hosts",
			beego.NSInclude(
				&ca.HostController{},
			),
		),
		beego.NSNamespace("/v1/asset/k8s",
			beego.NSInclude(
				&ca.K8SController{},
			),
		),
		beego.NSNamespace("/v1/securitypolicy/bmts",
			beego.NSInclude(
				&cs.BMTController{},
			),
		),
		beego.NSNamespace("/v1/securitylog",
			beego.NSInclude(
				&csl.IntrudeDetectLogController{},
				&csl.BenchMarkLogController{},
			),
		),
		beego.NSNamespace("/v1/merticinfo",
			beego.NSInclude(
				&controllers.MetricController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		beego.NSNamespace("/accounts",
			beego.NSInclude(
				&caccounts.AccountsController{},
			),
		),
		beego.NSNamespace("/groups",
			beego.NSInclude(
				&cgroups.GroupsController{},
			),
		),
		beego.NSNamespace("/v1/system",
			beego.NSInclude(
				&css.K8sController{},
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
	beego.InsertFilter("/api/v1/*", beego.BeforeRouter, isLogin)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.AddNamespace(ns)
}
