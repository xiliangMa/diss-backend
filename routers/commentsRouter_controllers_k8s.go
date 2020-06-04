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
			Method:           "AddCluster",
			Router:           `/add`,
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

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:PodController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/k8s:PodController"],
		beego.ControllerComments{
			Method:           "GetPodsList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
