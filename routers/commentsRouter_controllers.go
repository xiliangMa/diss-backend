package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers:AuthController"],
		beego.ControllerComments{
			Method:           "Authorize",
			Router:           `/Authorization`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers:AuthController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
