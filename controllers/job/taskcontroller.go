package job

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	taskservice "github.com/xiliangMa/diss-backend/service/task"
)

// Task 接口
type TaskController struct {
	beego.Controller
}

// @Title GetTaskList
// @Description Get Task List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.Task false "任务"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *TaskController) GetTaskList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	task := new(models.Task)
	json.Unmarshal(this.Ctx.Input.RequestBody, &task)
	this.Data["json"] = task.List(from, limit)
	this.ServeJSON(false)

}

// @Title DeleteTask
// @Description Delete Task
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id query string "" true "id"
// @Success 200 {object} models.Result
// @router / [delete]
func (this *TaskController) DeleteTask() {
	id := this.GetString("id")
	Task := new(models.Task)
	Task.Id = id
	for _, task := range Task.GetTaskList() {
		taskService := new(taskservice.TaskService)
		taskService.Task = task
		taskService.RemoveTask()
	}

	this.Data["json"] = models.Result{Code: http.StatusOK}
	this.ServeJSON(false)
}

// @Title GetTaskLogList
// @Description Get TaskLog List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.TaskLog false "任务调度"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /logs [post]
func (this *TaskController) GetTaskLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	taskLog := new(models.TaskLog)
	json.Unmarshal(this.Ctx.Input.RequestBody, &taskLog)
	this.Data["json"] = taskLog.List(from, limit)
	this.ServeJSON(false)

}
