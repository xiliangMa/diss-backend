package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"],
		beego.ControllerComments{
			Method:           "GetHostBnechMarkSummaryStatistics",
			Router:           `/:hostId/bms/host`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"],
		beego.ControllerComments{
			Method:           "GetAssetStatistics",
			Router:           `/asset`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"],
		beego.ControllerComments{
			Method:           "GetBnechMarkProportionStatistics",
			Router:           `/bmp`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"],
		beego.ControllerComments{
			Method:           "GetBnechMarkSummaryStatistics",
			Router:           `/bms`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/statistics:StatisticsController"],
		beego.ControllerComments{
			Method:           "GetIntrudeDetectLogStatistics",
			Router:           `/idl`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
