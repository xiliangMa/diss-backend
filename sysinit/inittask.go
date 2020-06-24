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
		// 启动 k8s 同步任务
		k8sSyncTaskId, _ := uuid.NewV4()
		syncSpec := beego.AppConfig.String("k8s::SyncSpec")
		k8STaskHandler := task.K8STaskHandler{}
		logs.Info("Start K8S Sync TaskHandler SyncSpec: %s", syncSpec)
		if err := th.AddByFunc(k8sSyncTaskId.String(), syncSpec, k8STaskHandler.SyncAll); err != nil {
			logs.Error("Start K8S Sync TaskHandler fail, err: %s", err)
			os.Exit(-1)
		} else {
			// to do 添加任务到数据库管理
		}

		// agent 系统检查任务
		systemCheckTaskId, _ := uuid.NewV4()
		systemCheckSpec := beego.AppConfig.String("system::SystemCheckSpec")
		logs.Info("Start System Check TaskHandler, SystemCheckSpec: %s", systemCheckSpec)
		if err := th.AddByFunc(systemCheckTaskId.String(), systemCheckSpec, new(task.SystemCheckHandler).SystemCheck); err != nil {
			logs.Error("Start System Check TaskHandler fail, err: %s", err)
			os.Exit(-1)
		}
		th.Start()
		select {}
	}()
}
