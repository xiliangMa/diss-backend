package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"],
		beego.ControllerComments{
			Method:           "GetImageMetadata",
			Router:           `/images/:imageDigest/metadata`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"],
		beego.ControllerComments{
			Method:           "GetImageContent",
			Router:           `/images/:imageId/content`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"],
		beego.ControllerComments{
			Method:           "GetImageSensitiveInfo",
			Router:           `/images/:imageId/sensitiveinfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"],
		beego.ControllerComments{
			Method:           "GetImageVuln",
			Router:           `/images/:imageId/vulns`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"] = append(beego.GlobalControllerRouter["github.com/xiliangMa/diss-backend/controllers/anchoreengine:AnchoreImageController"],
		beego.ControllerComments{
			Method:           "GetImageVulnStatistics",
			Router:           `/images/:imageId/vulnsstatistics`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
