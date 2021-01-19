// @APIVersion 1.0.0
// @Title DISS API
// @Description DISS API
// @Schemes https, http
package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/xiliangMa/diss-backend/controllers"
	caccounts "github.com/xiliangMa/diss-backend/controllers/accounts"
	"github.com/xiliangMa/diss-backend/controllers/anchoreengine"
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
	ns := web.NewNamespace("/api",
		web.NSNamespace("/v1/users",
			web.NSInclude(
				&cbase.UserController{},
			),
		),
		web.NSNamespace("/v1/statistics",
			web.NSInclude(
				&cstatistics.StatisticsController{},
				&cstatistics.PackageStatisticsController{},
			),
		),
		web.NSNamespace("/v1/asset/images",
			web.NSInclude(
				&casset.ImageController{},
			),
		),
		web.NSNamespace("/v1/asset/hosts",
			web.NSInclude(
				&casset.HostController{},
			),
		),
		web.NSNamespace("/v1/asset/k8s",
			web.NSInclude(
				&casset.K8SController{},
			),
		),
		web.NSNamespace("/v1/securitypolicy/systmps",
			web.NSInclude(
				&csecuritypolicy.SystemTemplateController{},
			),
		),
		web.NSNamespace("/v1/securitypolicy/systmpgroups",
			web.NSInclude(
				&csecuritypolicy.SystemTemplateGroupController{},
			),
		),
		web.NSNamespace("/v1/securitylog",
			web.NSInclude(
				&csecuritylog.IntrudeDetectLogController{},
				&csecuritylog.BenchMarkLogController{},
				&csecuritylog.VirusLogController{},
				&csecuritylog.VulnerabilitiesLogController{},
				&csecuritylog.WarningInfoController{},
			),
		),
		web.NSNamespace("/v1/securityaudit",
			web.NSInclude(
				&csecurityaudit.CmdHistoryController{},
				&csecurityaudit.DockerEventController{},
			),
		),
		web.NSNamespace("/auth",
			web.NSInclude(
				&controllers.AuthController{},
			),
		),
		web.NSNamespace("/v1/accounts",
			web.NSInclude(
				&caccounts.AccountsController{},
			),
		),
		web.NSNamespace("/v1/groups",
			web.NSInclude(
				&cbase.GroupsController{},
			),
		),
		web.NSNamespace("/v1/hosts",
			web.NSInclude(
				&cbase.HostController{},
			),
		),
		web.NSNamespace("/v1/containers",
			web.NSInclude(
				&cbase.ContainerController{},
			),
		),
		web.NSNamespace("/v1/images",
			web.NSInclude(
				&cbase.ImageController{},
			),
		),
		web.NSNamespace("/v1/k8s/clusters",
			web.NSInclude(
				&ck8s.ClusterController{},
			),
		),
		web.NSNamespace("/v1/k8s/namespaces",
			web.NSInclude(
				&ck8s.NSController{},
			),
		),
		web.NSNamespace("/v1/k8s/pods",
			web.NSInclude(
				&ck8s.PodController{},
			),
		),
		web.NSNamespace("/v1/k8s/services",
			web.NSInclude(
				&ck8s.ServiceController{},
			),
		),
		web.NSNamespace("/v1/k8s/deployment",
			web.NSInclude(
				&ck8s.DeploymentController{},
			),
		),
		web.NSNamespace("/v1/k8s/networkpolicy",
			web.NSInclude(
				&ck8s.NetworkPolicyController{},
			),
		),
		web.NSNamespace("/v1/system",
			web.NSInclude(
				&csystem.SystemController{},
				&csystem.IntegrationController{},
				&csystem.LicenseController{},
				&csystem.FeedsController{},
			),
		),
		web.NSNamespace("/v1/jobs",
			web.NSInclude(
				&cjob.JobController{},
			),
		),
		web.NSNamespace("/v1/tasks",
			web.NSInclude(
				&cjob.TaskController{},
			),
		),
		web.NSNamespace("/v1/securitycheck",
			web.NSInclude(
				&csecuritycheck.SecurityCheckController{},
			),
		),
		web.NSNamespace("/v1/anchore",
			web.NSInclude(
				&anchoreengine.AnchoreImageController{},
			),
		),
	)

	// add route for ws
	web.Router("/metrics", &ws.WSMetricController{}, "*:Metrics")

	var isLogin = func(ctx *context.Context) {
		jwtService := auth.JwtService{}
		if ctx.Request.Method != "OPTIONS" {
			_, code := jwtService.CheckToken(ctx.Input.Header("token"))
			if code != http.StatusOK {
				ctx.Redirect(http.StatusUnauthorized, "/swagger")
			}
		}

	}
	web.InsertFilter("/api/v1/*", web.BeforeRouter, isLogin)

	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	web.AddNamespace(ns)
}
