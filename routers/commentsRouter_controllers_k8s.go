package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"],
		beego.ControllerComments{
			Method:           "GetClusters",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"],
		beego.ControllerComments{
			Method:           "UpdateCluster",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"],
		beego.ControllerComments{
			Method:           "DeleteCluster",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"],
		beego.ControllerComments{
			Method:           "SyncCluster",
			Router:           `/:id/sync`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"],
		beego.ControllerComments{
			Method:           "AddCluster",
			Router:           `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ClusterController"],
		beego.ControllerComments{
			Method:           "ClusterSecurityCheck",
			Router:           `/securitycheck`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:DeploymentController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:DeploymentController"],
		beego.ControllerComments{
			Method:           "GetDeploymentList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"],
		beego.ControllerComments{
			Method:           "GetNameSpaceList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"],
		beego.ControllerComments{
			Method:           "UpdateNameSpace",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"],
		beego.ControllerComments{
			Method:           "BindAccount",
			Router:           `/:nsId/bindaccount`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NSController"],
		beego.ControllerComments{
			Method:           "UnBindAccount",
			Router:           `/:nsId/unbindaccount`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"],
		beego.ControllerComments{
			Method:           "GetNetworkPolicysList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"],
		beego.ControllerComments{
			Method:           "UpdateNetworkPolicy",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"],
		beego.ControllerComments{
			Method:           "DeleteNetworkPolicy",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:NetworkPolicyController"],
		beego.ControllerComments{
			Method:           "AddNetworkPolicy",
			Router:           `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:PodController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:PodController"],
		beego.ControllerComments{
			Method:           "GetPodsList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ServiceController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:ServiceController"],
		beego.ControllerComments{
			Method:           "GetServicesList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
