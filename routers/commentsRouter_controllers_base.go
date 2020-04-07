package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerCmdHistoryController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerCmdHistoryController"],
		beego.ControllerComments{
			Method:           "GetContainerCmdHistorys",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

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

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ContainerController"],
		beego.ControllerComments{
			Method:           "GetContainerInfo",
			Router:           `/:containerId/info`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostCmdHistoryController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostCmdHistoryController"],
		beego.ControllerComments{
			Method:           "GetHostCmdHistoryList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
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

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostController"],
		beego.ControllerComments{
			Method:           "UpdateHost",
			Router:           `/:hostId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:HostController"],
		beego.ControllerComments{
			Method:           "GetHostPsList",
			Router:           `/:hostId/ps`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ImageController"],
		beego.ControllerComments{
			Method:           "GetImagesList",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ImageController"],
		beego.ControllerComments{
			Method:           "DeleteImage",
			Router:           `/:imageId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/base:ImageController"],
		beego.ControllerComments{
			Method:           "GetImageInfo",
			Router:           `/:imageId/info`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
