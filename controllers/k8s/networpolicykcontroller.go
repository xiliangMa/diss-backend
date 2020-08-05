package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

// NetworkPolicy 接口
type NetworkPolicyController struct {
	beego.Controller
}

// @Title GetNetworkPolicy
// @Description Get NetworkPolicy List
// @Param token header string true "authToken"
// @Param body body models.NetworkPolicy false "网络策略"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *NetworkPolicyController) GetNetworkPolicysList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	NetworkPolicy := new(models.NetworkPolicy)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NetworkPolicy)
	this.Data["json"] = NetworkPolicy.List(from, limit)
	this.ServeJSON(false)
}

// @Title AddNetworkPolicy
// @Description Add NetworkPolicy
// @Param token header string true "authToken"
// @Param body body models.NetworkPolicy true "网络策略"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *NetworkPolicyController) AddNetworkPolicy() {
	NetworkPolicy := new(models.NetworkPolicy)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NetworkPolicy)
	this.Data["json"] = NetworkPolicy.Add()
	this.ServeJSON(false)
}

// @Title UpdateNetworkPolicy
// @Description Update NetworkPolicy
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Param body body models.NetworkPolicy true "网络策略"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *NetworkPolicyController) UpdateNetworkPolicy() {
	id := this.GetString(":id")
	NetworkPolicy := new(models.NetworkPolicy)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NetworkPolicy)
	NetworkPolicy.Id = id
	this.Data["json"] = NetworkPolicy.Update()
	this.ServeJSON(false)
}

// @Title DeleteNetworkPolicy
// @Description Delete NetworkPolicy
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *NetworkPolicyController) DeleteNetworkPolicy() {
	id := this.GetString(":id")
	NetworkPolicy := new(models.NetworkPolicy)
	NetworkPolicy.Id = id
	NetworkPolicy.Delete()
	this.Data["json"] = models.Result{Code: http.StatusOK}
	this.ServeJSON(false)
}
