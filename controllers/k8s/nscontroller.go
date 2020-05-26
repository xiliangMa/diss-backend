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

// @Title UpdateNameSpaces
// @Description Update NameSpaces
// @Param token header string true "authToken"
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
