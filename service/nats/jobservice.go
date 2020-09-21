package nats

import (
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type JobService struct {
	Task *models.Task
}

func (this *JobService) SetJobDeactiveStatus() models.Result {
	result := models.Result{Code: http.StatusOK}

	taskParam := new(models.Task)
	taskParam.IsOne = true
	taskParam.Type = models.Job_Type_Periodic
	taskParam.Status = models.Task_Status_Unavailable
	unavailableTask := taskParam.List(0, 1)
	if unavailableTask.Data == nil {
		job := this.Task.Job
		jobObj := job.Get()
		jobObj.Status = models.Job_Status_Deactived
		result = jobObj.Update()
	}

	return result
}
