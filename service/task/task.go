package task

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"sync"
)

var (
	secondParser = cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
)

type TaskHandler struct {
	c     *cron.Cron
	ids   map[string]cron.EntryID
	mutex sync.Mutex
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		c:   cron.New(cron.WithParser(secondParser), cron.WithChain()),
		ids: make(map[string]cron.EntryID),
	}
}

// Start start the TaskHandler engine
func (this *TaskHandler) Start() {
	this.c.Start()
}

// Stop stop the TaskHandler engine
func (this *TaskHandler) Stop() {
	this.c.Stop()
}

// DelByID remove one TaskHandler task
func (this *TaskHandler) DelByID(id string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	eid, ok := this.ids[id]
	if !ok {
		return
	}
	this.c.Remove(eid)
	delete(this.ids, id)
}

// AddByID add one TaskHandler task
// id is unique
// spec is the TaskHandler expression
func (this *TaskHandler) AddByID(id string, spec string, cmd cron.Job) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.ids[id]; ok {
		return errors.Errorf("TaskHandler id exists")
	}
	eid, err := this.c.AddJob(spec, cmd)
	if err != nil {
		return err
	}
	this.ids[id] = eid
	return nil
}

// AddByFunc add function as TaskHandler task
func (this *TaskHandler) AddByFunc(id string, spec string, f func()) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.ids[id]; ok {
		return errors.Errorf("TaskHandler id exists")
	}
	eid, err := this.c.AddFunc(spec, f)
	if err != nil {
		return err
	}
	this.ids[id] = eid
	return nil
}

// IsExists check the TaskHandler task whether existed with c id
func (this *TaskHandler) IsExists(jid string) bool {
	_, exist := this.ids[jid]
	return exist
}
