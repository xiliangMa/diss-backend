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

}
