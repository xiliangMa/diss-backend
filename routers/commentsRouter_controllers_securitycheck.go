package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitycheck:SecurityCheckController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitycheck:SecurityCheckController"],
		beego.ControllerComments{
			Method:           "SecurityCheck",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitycheck:SecurityCheckController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitycheck:SecurityCheckController"],
		beego.ControllerComments{
			Method:           "SecurityCheck2",
			Router:           `/v2`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
