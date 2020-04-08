package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:BenchMarkLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:BenchMarkLogController"],
        beego.ControllerComments{
            Method: "GetBenchMarkLogList",
            Router: `/bmls`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"],
        beego.ControllerComments{
            Method: "GetIntrudeDetectLogList",
            Router: `/idls`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"],
        beego.ControllerComments{
            Method: "GetIntrudeDetectLogInfo",
            Router: `/intrudedetect/:hostId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
