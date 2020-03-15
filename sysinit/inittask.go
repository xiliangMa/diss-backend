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
		// 启动 k8s 任务
		uid, _ := uuid.NewV4()
		if err := th.AddByFunc(uid.String(), beego.AppConfig.String("k8s::SyncSpec"), task.SyncAll); err != nil {
			logs.Error("error to add TaskHandler task: %s", err)
			os.Exit(-1)
		} else {
			// to do 添加任务到数据库管理
		}
		th.Start()
		select {}
	}()
}
