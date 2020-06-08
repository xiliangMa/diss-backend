package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostConfigList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostInfoList",
			Router:           `/:hostId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostCmdHistoryList",
			Router:           `/:hostId/cmdhistory`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostBenchMarkLogList",
			Router:           `/:hostId/hostbmls`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostBenchMarkLogInfo",
			Router:           `/:hostId/hostbmls/:bmlId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostImagesList",
			Router:           `/:hostId/images`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostImageInfo",
			Router:           `/:hostId/images/:imageId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostPsList",
			Router:           `/:hostId/ps`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostContainerConfigList",
			Router:           `/:hostName/containers`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostContainerInfoList",
			Router:           `/:hostName/containers/:containerId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:HostController"],
		beego.ControllerComments{
			Method:           "GetHostPodList",
			Router:           `/:hostName/pods`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:ImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:ImageController"],
		beego.ControllerComments{
			Method:           "GetContainersList",
			Router:           `/containers`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetContainerImageInfo",
			Router:           `/:hostName/imageinfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetClusters",
			Router:           `/clusters`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetNameSpaces",
			Router:           `/clusters/:clusterId/namespaces`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetContainerInfo",
			Router:           `/containers/:containerId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetContainerCmdHistorys",
			Router:           `/containers/:containerId/cmdhistorys`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetContainerPs",
			Router:           `/containers/:containerId/containerps`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "BindAccount",
			Router:           `/namespaces/:nsId/bindaccount`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "UnBindAccount",
			Router:           `/namespaces/:nsId/unbindaccount`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetPods",
			Router:           `/namespaces/:nsName/pods`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/asset:K8SController"],
		beego.ControllerComments{
			Method:           "GetContainerConfig",
			Router:           `/namespaces/:nsName/pods/:podId/containers`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
