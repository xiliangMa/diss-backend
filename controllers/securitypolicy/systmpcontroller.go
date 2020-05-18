package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// System Template object api list
type SystemTemplateController struct {
	beego.Controller
}

// @Title GetSystemTemplateList
// @Description Get System Template List
// @Param token header string true "authToken"
// @Param body body models.SystemTemplate false "安全策略"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SystemTemplateController) GetSystemTemplateLIst() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	systemTemplate := new(models.SystemTemplate)
	json.Unmarshal(this.Ctx.Input.RequestBody, &systemTemplate)
	this.Data["json"] = systemTemplate.List(from, limit)
	this.ServeJSON(false)
}
