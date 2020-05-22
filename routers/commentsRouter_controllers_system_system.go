package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
		beego.ControllerComments{
			Method:           "GetLicenseData",
			Router:           `/system/license`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
		beego.ControllerComments{
			Method:           "AddLicenseFile",
			Router:           `/system/license/import`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
		beego.ControllerComments{
			Method:           "AddLogConfig",
			Router:           `/system/logconfig`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
		beego.ControllerComments{
			Method:           "UpdateLogConfig",
			Router:           `/system/logconfig`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
		beego.ControllerComments{
			Method:           "GetLogConfig",
			Router:           `/system/logconfig`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
		beego.ControllerComments{
			Method:           "UploadLogo",
			Router:           `/system/logo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
		beego.ControllerComments{
			Method:           "CheckLogoIsExist",
			Router:           `/system/logo/isexist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
