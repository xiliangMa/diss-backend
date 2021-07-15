package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
)

type CommonService struct {
	ImageConfig *models.ImageConfig
	Token       string
}

var scheme = regexp.MustCompile("(https|http)://([-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|])")

func (this *CommonService) AddDetail() {
	var ref name.Reference
	rs := scheme.FindStringSubmatch(this.ImageConfig.Registry.Url)
	ref, _ = name.ParseReference(this.ImageConfig.Name)

	if rs != nil && this.ImageConfig.Registry.Type != models.Registry_Type_DockerHub {
		if this.ImageConfig.Namespaces != "" && this.ImageConfig.Registry.Type != models.Registry_Type_AwsECR {
			this.ImageConfig.Name = strings.Replace(rs[2], "/", "", 1) + "/" + this.ImageConfig.Namespaces + "/" + this.ImageConfig.Name
		} else {
			this.ImageConfig.Name = strings.Replace(rs[2], "/", "", 1) + "/" + this.ImageConfig.Name
		}
		ref, _ = name.ParseReference(this.ImageConfig.Name)
		if rs[1] == "http" {
			ref, _ = name.ParseReference(this.ImageConfig.Name, name.Insecure)
		}
	}
	cs := CommonService{ImageConfig: this.ImageConfig}
	task := cs.AddTask()
	msg := ""

	want := authn.AuthConfig{Username: this.ImageConfig.Registry.User, Password: this.ImageConfig.Registry.Pwd}
	if this.Token != "" {
		want = authn.AuthConfig{Auth: this.Token}
	}
	img, err := remote.Image(ref, remote.WithAuth(authn.FromConfig(want)))
	if err != nil {
		task.Status = models.Task_Status_Failed
		msg = err.Error()
	}

	hash, _ := img.ConfigName()
	cf, _ := img.ConfigFile()
	digest, _ := img.Digest()
	layer, _ := img.Layers()

	var sum int64
	for _, l := range layer {
		size, _ := l.Size()
		sum += size
	}
	this.ImageConfig.Id = ""
	this.ImageConfig.ImageId = hash.String()
	if ic := this.ImageConfig.Get(); ic == nil {
		this.ImageConfig.Size = utils.FormatFileSize(sum)
		this.ImageConfig.CreateTime = cf.Created.UnixNano()
		this.ImageConfig.Add()

		imageDetail := models.ImageDetail{}
		imageDetail.ImageId = hash.String()
		imageDetail.Name = this.ImageConfig.Name
		imageDetail.ImageConfigId = this.ImageConfig.Id
		imageDetail.Layers = len(layer)
		imageDetail.RepoDigests = digest.String()
		imageDetail.Size = this.ImageConfig.Size
		imageDetail.CreateTime = cf.Created.UnixNano()

		var buffer bytes.Buffer
		var trimnop = regexp.MustCompile(`^/bin/sh\s+-c\s+#\(nop\)\s+`)
		var trimrun = regexp.MustCompile(`^(RUN\s+){0,1}/bin/sh\s+-c\s+`)
		for _, h := range cf.History {
			tmpstr := trimnop.ReplaceAllString(h.CreatedBy, "")
			tmpstrs := trimrun.ReplaceAllString(tmpstr, "RUN ")
			buffer.WriteString(tmpstrs + "\n")
		}
		imageDetail.Dockerfile = strings.TrimSpace(buffer.String())
		imageDetail.Add()

		task.Status = models.Task_Status_Finished
		task.RunCount = 1
	} else {
		task.Status = models.Task_Status_Failed
		msg = "镜像已存在"
	}
	task.Update()
	taskRawInfo, _ := json.Marshal(task)
	if msg == "" {
		msg = fmt.Sprintf("更新任务成功, 状态: %s >>> 镜像名: %s, 任务ID: %s <<<", "完成", this.ImageConfig.Name, task.Id)
	} else {
		msg = fmt.Sprintf("更新任务失败, 状态: %s >>> 镜像名: %s, 任务ID: %s 失败原因: %s <<<", "失败", this.ImageConfig.Name, task.Id, msg)
	}
	taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Info}
	taskLog.Add()
}

func (this *CommonService) AddTask() *models.Task {
	task := models.Task{}
	uid, _ := uuid.NewV4()
	task.Id = uid.String()
	taskpre := "用户任务-"
	task.Name = taskpre + task.Id
	task.Type = models.Job_Type_Once
	task.Description = taskpre + "仓库镜像批量导入"
	task.Image = this.ImageConfig
	task.Batch = time.Now().UnixNano() / 1e3
	task.Status = models.Task_Status_Created
	t := task.Add()

	taskLog := models.TaskLog{}
	taskRawInfo, _ := json.Marshal(task)
	taskLog.Task = string(taskRawInfo)
	taskLog.Level = models.Log_level_Info
	msg := fmt.Sprintf("创建任务成功, 状态: 已创建,  批次: %v, 任务ID: %s", task.Batch, task.Id)
	taskLog.RawLog = msg
	taskLog.Add()

	if t.Data != nil {
		return t.Data.(*models.Task)
	}
	return nil

}
