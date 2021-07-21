package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 命名空间接口
type NSController struct {
	beego.Controller
}

// @Title GetNameSpaces
// @Description Get NameSpace List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.NameSpace false "命名空间"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *NSController) GetNameSpaceList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	ns := new(models.NameSpace)
	json.Unmarshal(this.Ctx.Input.RequestBody, &ns)
	this.Data["json"] = ns.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateNameSpace
// @Description Update NameSpace
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "Id"
// @Param body body models.NameSpace true "命名空间"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *NSController) UpdateNameSpace() {
	id := this.GetString(":id")
	nameSpace := new(models.NameSpace)
	json.Unmarshal(this.Ctx.Input.RequestBody, &nameSpace)
	nameSpace.Id = id
	this.Data["json"] = nameSpace.Update()
	this.ServeJSON(false)
}

// @Title BindAccount
// @Description BindAccount（绑定租户）
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param nsId path string "" true "nsId"
// @Param body body models.NameSpace true "命名空间"
// @Success 200 {object} models.Result
// @router /:nsId/bindaccount [put]
func (this *NSController) BindAccount() {
	nsId := this.GetString(":nsId")
	NS := new(models.NameSpace)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NS)
	NS.Id = nsId
	this.Data["json"] = NS.BindAccount()
	this.ServeJSON(false)
}

// @Title UnBindAccount
// @Description UnBindAccount（解除绑定）
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param nsId path string "" true "nsId"
// @Param body body models.NameSpace true "命名空间"
// @Success 200 {object} models.Result
// @router /:nsId/unbindaccount [delete]
func (this *NSController) UnBindAccount() {
	nsId := this.GetString(":nsId")
	NS := new(models.NameSpace)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NS)
	NS.Id = nsId
	this.Data["json"] = NS.UnBindAccount()
	this.ServeJSON(false)
}
