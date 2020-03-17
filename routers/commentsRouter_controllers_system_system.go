package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:K8sController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:K8sController"],
		beego.ControllerComments{
			Method:           "UploadK8sFile",
			Router:           `/system/k8s/upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
