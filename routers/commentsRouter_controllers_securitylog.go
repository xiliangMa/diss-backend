package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:BenchMarkLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:BenchMarkLogController"],
		beego.ControllerComments{
			Method:           "GetBenchMarkLogList",
			Router:           `/bmls`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"],
		beego.ControllerComments{
			Method:           "GetIntrudeDetectLogList",
			Router:           `/idls`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:IntrudeDetectLogController"],
		beego.ControllerComments{
			Method:           "GetIntrudeDetectLogInfo",
			Router:           `/intrudedetect/:hostId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VirusLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VirusLogController"],
		beego.ControllerComments{
			Method:           "GetHostOrContainerVirusLogList",
			Router:           `/virus/hostorcontainer`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VirusLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VirusLogController"],
		beego.ControllerComments{
			Method:           "GetVirusLogList",
			Router:           `/virus/image`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VulnerabilitiesLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VulnerabilitiesLogController"],
		beego.ControllerComments{
			Method:           "GetImageVulnerabilitiesLogList",
			Router:           `/vulnerabilities/image`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VulnerabilitiesLogController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:VulnerabilitiesLogController"],
		beego.ControllerComments{
			Method:           "GetImageVulnerabilityInfo",
			Router:           `/vulnerabilities/image/info`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:WarningInfoController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:WarningInfoController"],
		beego.ControllerComments{
			Method:           "GetWarningInfoList",
			Router:           `/warninginfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:WarningInfoController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:WarningInfoController"],
		beego.ControllerComments{
			Method:           "UpdateWarningInfo",
			Router:           `/warninginfo/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:WarningInfoController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitylog:WarningInfoController"],
		beego.ControllerComments{
			Method:           "AddClientSub_Image_Safe",
			Router:           `/warninginfo/addsub_image/:libname`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
