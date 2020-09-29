package sysinit

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/service/task"
	"os"
)

func InitTask() {
	go func() {
		th := task.NewTaskHandler()
		InitSystemCheckTask(th)
		InitClusterStatusCheckTask(th)
		th.Start()
		select {}
	}()
}

/**
 * 安全检查任务
 */
func InitSystemCheckTask(th *task.TaskHandler) {
	// agent 系统检查任务
	systemCheckTaskId, _ := uuid.NewV4()
	systemCheckSpec := beego.AppConfig.String("system::SystemCheckSpec")
	logs.Info("Start System Check TaskHandler, SystemCheckSpec: %s", systemCheckSpec)
	if err := th.AddByFunc(systemCheckTaskId.String(), systemCheckSpec, new(task.SystemCheckHandler).SystemCheck); err != nil {
		logs.Error("Start System Check TaskHandler fail, err: %s", err)
		os.Exit(-1)
	}
}

/**
 * 集群状态检查
 */
func InitClusterStatusCheckTask(th *task.TaskHandler) {
	clusterStatusCheckTaskId, _ := uuid.NewV4()
	clusterStatusCheckSpec := beego.AppConfig.String("system::ClusterStatusCheckSpec")
	logs.Info("Start Cluster Check TaskHandler, clusterStatusCheckSpec: %s", clusterStatusCheckSpec)
	if err := th.AddByFunc(clusterStatusCheckTaskId.String(), clusterStatusCheckSpec, new(task.K8sTaskHandler).CheckClusterStatusTask); err != nil {
		logs.Error("Start Cluster Check TaskHandler fail, err: %s", err)
		os.Exit(-1)
	}
}
