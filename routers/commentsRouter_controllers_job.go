package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"],
        beego.ControllerComments{
            Method: "GetTaskList",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
