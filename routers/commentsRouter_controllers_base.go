package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerController"],
		beego.ControllerComments{
			Method:           "GetContainersList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerController"],
		beego.ControllerComments{
			Method:           "DeleteContainer",
			Router:           `/:containerId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostController"],
		beego.ControllerComments{
			Method:           "HostList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
