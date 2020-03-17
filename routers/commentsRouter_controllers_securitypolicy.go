package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:BMTController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:BMTController"],
		beego.ControllerComments{
			Method:           "GetBMTList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
