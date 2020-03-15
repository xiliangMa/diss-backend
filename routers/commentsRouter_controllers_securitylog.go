package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"],
        beego.ControllerComments{
            Method: "GetIntrudeDetectLogInfo",
            Router: `/intrudedetect/:hostId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
