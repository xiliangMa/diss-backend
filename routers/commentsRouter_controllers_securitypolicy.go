package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateController"],
		beego.ControllerComments{
			Method:           "GetSystemTemplateLIst",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
