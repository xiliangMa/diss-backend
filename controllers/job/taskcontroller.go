package job

import (
	"encoding/json"
	"github.com/astaxie/beego"
	mjob "github.com/xiliangMa/diss-backend/models/job"
)

// Task 接口
type TaskController struct {
	beego.Controller
}

// @Title GetTaskList
// @Description Get Task List
// @Param token header string true "authToken"
// @Param body body job.Task false "任务"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *TaskController) GetTaskList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	task := new(mjob.Task)
	json.Unmarshal(this.Ctx.Input.RequestBody, &task)
	this.Data["json"] = task.List(from, limit)
	this.ServeJSON(false)

}

// @Title DeleteTask
// @Description Delete Task
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *TaskController) DeleteTask() {
	id := this.GetString(":id")
	task := new(mjob.Task)
	task.Id = id
	this.Data["json"] = task.Delete()
	this.ServeJSON(false)

}
