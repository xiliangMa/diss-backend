package job

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/ws"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

// Task 接口
type TaskController struct {
	beego.Controller
}

// @Title GetTaskList
// @Description Get Task List
// @Param token header string true "authToken"
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
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *TaskController) DeleteTask() {
	id := this.GetString(":id")
	task := new(models.Task)
	task.Id = id
	result := task.List(0, 0)
	if result.Data != nil {
		data := result.Data.(map[string]interface{})
		if result.Code == http.StatusOK && data["total"] != 0 {
			//向agent下发删除任务指令
			deleteTaskList := data["items"]
			deleteTask := deleteTaskList.([]*models.Task)[0]
			WSDeliverService := ws.WSDeliverService{Hub: models.WSHub, DelTask: deleteTask}
			err := WSDeliverService.DeleteTask()
			if err == nil {
				// 更新数据库状态 设置为删除锁定
				deleteTask.Status = models.Task_Status_Removing
				deleteTask.Update()
			} else {
				result.Code = utils.DeleteTaskErr
				result.Message = "DeleteTaskErr"
				result.Data = nil
				msg := fmt.Sprintf("Delet Task fail, Id: %s, err: %s", task.Id, result.Message)
				taskLog := models.TaskLog{RawLog: msg, Task: task, Level: models.Log_level_Error}
				taskLog.Add()
			}
		}
	}
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title GetTaskLogList
// @Description Get TaskLog List
// @Param token header string true "authToken"
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
