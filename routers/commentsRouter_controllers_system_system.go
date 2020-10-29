package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:FeedsController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:FeedsController"],
        beego.ControllerComments{
            Method: "GetFeedList",
            Router: `/feeds`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
        beego.ControllerComments{
            Method: "AddLogConfig",
            Router: `/system/logconfig`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
        beego.ControllerComments{
            Method: "UpdateLogConfig",
            Router: `/system/logconfig`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:IntegrationController"],
        beego.ControllerComments{
            Method: "GetLogConfig",
            Router: `/system/logconfig`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "GetFeatureCode",
            Router: `/system/featurecode`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "VerifyFeatureCode",
            Router: `/system/featurecode/verify`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "GetLicense",
            Router: `/system/license`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "GetLicensedHostCount",
            Router: `/system/license/hostcount`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "SetHostLicense",
            Router: `/system/license/hostlicense/:hostId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "AddLicenseFile",
            Router: `/system/license/import`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:LicenseController"],
        beego.ControllerComments{
            Method: "GetLicenseHistory",
            Router: `/system/licensehistory`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
        beego.ControllerComments{
            Method: "UploadLogo",
            Router: `/system/logo`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
        beego.ControllerComments{
            Method: "CheckLogoIsExist",
            Router: `/system/logo/isexist`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
        beego.ControllerComments{
            Method: "AddSysConfig",
            Router: `/system/sysconfig`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
        beego.ControllerComments{
            Method: "GetSysConfigs",
            Router: `/system/sysconfig`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/system/system:SystemController"],
        beego.ControllerComments{
            Method: "UpdateSysConfig",
            Router: `/system/sysconfig/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
