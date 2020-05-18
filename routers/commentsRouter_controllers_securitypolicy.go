package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateController"],
        beego.ControllerComments{
            Method: "GetSystemTemplateLIst",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateGroupController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateGroupController"],
        beego.ControllerComments{
            Method: "GetSystemTemplateGroupLIst",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateGroupController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/securitypolicy:SystemTemplateGroupController"],
        beego.ControllerComments{
            Method: "AddSystemTemplateGroup",
            Router: `/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
