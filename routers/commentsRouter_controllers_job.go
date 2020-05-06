package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:JobController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:JobController"],
        beego.ControllerComments{
            Method: "GetJobList",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:JobController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:JobController"],
        beego.ControllerComments{
            Method: "DeleteJob",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:JobController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:JobController"],
        beego.ControllerComments{
            Method: "AddJob",
            Router: `/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"],
        beego.ControllerComments{
            Method: "GetTaskList",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"],
        beego.ControllerComments{
            Method: "DeleteTask",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/job:TaskController"],
        beego.ControllerComments{
            Method: "GetTaskLogList",
            Router: `/logs`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
