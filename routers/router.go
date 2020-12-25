// @APIVersion 1.0.0
// @Title DISS API
// @Description DISS API
// @Schemes https, http
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/xiliangMa/diss-backend/controllers"
	caccounts "github.com/xiliangMa/diss-backend/controllers/accounts"
	casset "github.com/xiliangMa/diss-backend/controllers/asset"
	cbase "github.com/xiliangMa/diss-backend/controllers/base"
	cjob "github.com/xiliangMa/diss-backend/controllers/job"
	ck8s "github.com/xiliangMa/diss-backend/controllers/k8s"
	csecurityaudit "github.com/xiliangMa/diss-backend/controllers/securityaudit"
	csecuritycheck "github.com/xiliangMa/diss-backend/controllers/securitycheck"
	csecuritylog "github.com/xiliangMa/diss-backend/controllers/securitylog"
	csecuritypolicy "github.com/xiliangMa/diss-backend/controllers/securitypolicy"
	cstatistics "github.com/xiliangMa/diss-backend/controllers/statistics"
	csystem "github.com/xiliangMa/diss-backend/controllers/system/system"
	"github.com/xiliangMa/diss-backend/service/auth"

	ws "github.com/xiliangMa/diss-backend/controllers/ws"
	"net/http"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1/users",
			beego.NSInclude(
				&cbase.UserController{},
			),
		),
		beego.NSNamespace("/v1/statistics",
			beego.NSInclude(
				&cstatistics.StatisticsController{},
				&cstatistics.PackageStatisticsController{},
			),
		),
		beego.NSNamespace("/v1/asset/images",
			beego.NSInclude(
				&casset.ImageController{},
			),
		),
		beego.NSNamespace("/v1/asset/hosts",
			beego.NSInclude(
				&casset.HostController{},
			),
		),
		beego.NSNamespace("/v1/asset/k8s",
			beego.NSInclude(
				&casset.K8SController{},
			),
		),
		beego.NSNamespace("/v1/securitypolicy/systmps",
			beego.NSInclude(
				&csecuritypolicy.SystemTemplateController{},
			),
		),
		beego.NSNamespace("/v1/securitypolicy/systmpgroups",
			beego.NSInclude(
				&csecuritypolicy.SystemTemplateGroupController{},
			),
		),
		beego.NSNamespace("/v1/securitylog",
			beego.NSInclude(
				&csecuritylog.IntrudeDetectLogController{},
				&csecuritylog.BenchMarkLogController{},
				&csecuritylog.VirusLogController{},
				&csecuritylog.VulnerabilitiesLogController{},
				&csecuritylog.WarningInfoController{},
			),
		),
		beego.NSNamespace("/v1/securityaudit",
			beego.NSInclude(
				&csecurityaudit.CmdHistoryController{},
				&csecurityaudit.DockerEventController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		beego.NSNamespace("/v1/accounts",
			beego.NSInclude(
				&caccounts.AccountsController{},
			),
		),
		beego.NSNamespace("/v1/groups",
			beego.NSInclude(
				&cbase.GroupsController{},
			),
		),
		beego.NSNamespace("/v1/hosts",
			beego.NSInclude(
				&cbase.HostController{},
			),
		),
		beego.NSNamespace("/v1/containers",
			beego.NSInclude(
				&cbase.ContainerController{},
			),
		),
		beego.NSNamespace("/v1/images",
			beego.NSInclude(
				&cbase.ImageController{},
			),
		),
		beego.NSNamespace("/v1/k8s/clusters",
			beego.NSInclude(
				&ck8s.ClusterController{},
			),
		),
		beego.NSNamespace("/v1/k8s/namespaces",
			beego.NSInclude(
				&ck8s.NSController{},
			),
		),
		beego.NSNamespace("/v1/k8s/pods",
			beego.NSInclude(
				&ck8s.PodController{},
			),
		),
		beego.NSNamespace("/v1/k8s/services",
			beego.NSInclude(
				&ck8s.ServiceController{},
			),
		),
		beego.NSNamespace("/v1/k8s/deployment",
			beego.NSInclude(
				&ck8s.DeploymentController{},
			),
		),
		beego.NSNamespace("/v1/k8s/networkpolicy",
			beego.NSInclude(
				&ck8s.NetworkPolicyController{},
			),
		),
		beego.NSNamespace("/v1/system",
			beego.NSInclude(
				&csystem.SystemController{},
				&csystem.IntegrationController{},
				&csystem.LicenseController{},
				&csystem.FeedsController{},
			),
		),
		beego.NSNamespace("/v1/jobs",
			beego.NSInclude(
				&cjob.JobController{},
			),
		),
		beego.NSNamespace("/v1/tasks",
			beego.NSInclude(
				&cjob.TaskController{},
			),
		),
		beego.NSNamespace("/v1/securitycheck",
			beego.NSInclude(
				&csecuritycheck.SecurityCheckController{},
			),
		),
	)

	// add route for ws
	beego.Router("/metrics", &ws.WSMetricController{}, "*:Metrics")

	var isLogin = func(ctx *context.Context) {
		jwtService := auth.JwtService{}
		if ctx.Request.Method != "OPTIONS" {
			_, code := jwtService.CheckToken(ctx.Input.Header("token"))
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
