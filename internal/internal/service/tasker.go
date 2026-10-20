package service

import (
	"context"
	"time"

	"github.com/goexl/exception"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/goexl/schedule"
	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/core"
	"github.com/harluo/taskd/internal/internal/internal/db"
	"github.com/harluo/taskd/internal/internal/internal/get"
	"github.com/harluo/taskd/internal/internal/internal/model"
)

type Tasker struct {
	schedule db.Schedule
	task     db.Task

	scheduler *schedule.Scheduler
	runnable  *core.Runnable
	logger    log.Logger

	id string
}

func newTasker(tasker get.Tasker) *Tasker {
	return &Tasker{
		schedule: tasker.Schedule,
		task:     tasker.Task,

		scheduler: tasker.Scheduler,
		logger:    tasker.Logger,
		runnable:  tasker.Runnable,
	}
}

func (t *Tasker) Start(_ context.Context) (err error) {
	name := "任务执行器"
	fields := gox.Fields[any]{
		field.New("name", name),
	}
	if rid, ae := t.scheduler.Add(t).Duration(3 * time.Second).Name(name).Build().Apply(); nil != ae {
		err = ae
		t.logger.Error("添加任务出错", field.Error(ae), fields...)
	} else {
		t.id = rid
		t.logger.Info("添加任务成功", field.New("id", rid), fields...)
	}

	return
}

func (t *Tasker) Stop(_ context.Context) (err error) {
	t.scheduler.Remove().Id(t.id).Build().Apply()
	err = t.scheduler.Stop()

	return
}

func (t *Tasker) Running(id uint64, status task.Status, times uint64) (err error) {
	_task := new(model.Task)
	_task.Id = id
	_task.Status = status
	_task.Times = times
	if _, ue := t.task.Update(context.Background(), _task); nil != ue {
		err = ue
	}

	return
}

func (t *Tasker) Update(id uint64, status task.Status, runtime time.Time) (err error) {
	_task := new(model.Task)
	_task.Id = id
	_task.Status = status
	_task.Next = runtime
	if _, ue := t.task.Update(context.Background(), _task); nil != ue {
		err = ue
	}

	return
}

func (t *Tasker) Pop() task.Task {
	return t.runnable.Task()
}

func (t *Tasker) Archive(task task.Task) (err error) {
	_task := new(model.Task)
	_task.Id = task.Id()
	if exists, ge := t.task.Get(context.Background(), _task); nil != ge {
		err = ge
	} else if !exists {
		err = exception.New().Message("任务不存在").Field(field.New("task", task)).Build()
	} else if _, ae := t.task.Archive(context.Background(), _task); nil != ae {
		err = ae
	}

	return
}

func (t *Tasker) Failed(_ task.Task) (err error) {
	return
}

func (t *Tasker) Run() (err error) {
	if tasks, re := t.task.GetsRunnable(context.Background()); nil != re {
		err = re
	} else if 0 != len(*tasks) {
		all := *tasks
		t.runnable.Put(all[0], all[1:]...)
	}

	return
}
