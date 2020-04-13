package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	msecuritypolicy "github.com/xiliangMa/diss-backend/models/securitypolicy"
)

// System Template object api list
type SystemTemplateController struct {
	beego.Controller
}

// @Title GetSystemTemplateList
// @Description Get System Template List
// @Param token header string true "authToken"
// @Param body body securitypolicy.SystemTemplate false "系统模版"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SystemTemplateController) GetSystemTemplateLIst() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	systemTemplate := new(msecuritypolicy.SystemTemplate)
	json.Unmarshal(this.Ctx.Input.RequestBody, &systemTemplate)
	this.Data["json"] = systemTemplate.List(from, limit)
	this.ServeJSON(false)

}
