package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
		beego.ControllerComments{
			Method:           "UploadK8sFile",
			Router:           `/system/k8s/upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
		beego.ControllerComments{
			Method:           "UploadLogo",
			Router:           `/system/logo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
