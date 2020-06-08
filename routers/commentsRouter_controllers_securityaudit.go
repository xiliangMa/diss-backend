package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securityaudit:CmdHistoryController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securityaudit:CmdHistoryController"],
        beego.ControllerComments{
            Method: "GetCmdHistorys",
            Router: `/cmdhistorys`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
