package sysinit

import (
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/service/task"
	"os"
)

const (
	SyncK8sSpec = "0 */10 * * * ?"
)

func InitTask() {

	go func() {
		th := task.NewTaskHandler()
		// k8s
		uid, _ := uuid.NewV4()
		if err := th.AddByFunc(uid.String(), SyncK8sSpec, task.SyncAll); err != nil {
			logs.Error("error to add TaskHandler task: %s", err)
			os.Exit(-1)
		} else {
			// to do 添加任务记录
		}
		th.Start()
		select {}
	}()

}
