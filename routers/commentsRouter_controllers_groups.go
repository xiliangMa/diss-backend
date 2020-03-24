package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/groups:GroupsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/groups:GroupsController"],
        beego.ControllerComments{
            Method: "GetGroupsList",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/groups:GroupsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/groups:GroupsController"],
        beego.ControllerComments{
            Method: "GetContainersList",
            Router: `/containers`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
