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

func (t *Tasker) Add(required task.Schedule, optionals ...task.Schedule) (err error) {
	runtimes := make([]*model.Runtime, 0, 1+len(optionals))
	for _, _schedule := range append([]task.Schedule{required}, optionals...) {
		runtime := new(model.Runtime)
		runtime.Type = _schedule.Type()
		runtime.Next = _schedule.Next()
		runtime.Subtype = _schedule.Subtype()
		runtime.Target = _schedule.Target()
		runtime.Maximum = _schedule.Maximum()
		runtime.Timeout = _schedule.Timeout()
		runtime.Data = _schedule.Data()

		runtimes = append(runtimes, runtime)
	}
	if successes, ae := t.schedule.Add(runtimes[0], runtimes[1:]...); nil != ae {
		err = ae
	} else {
		t.runnable.Put((*successes)[0], (*successes)[1:]...)
	}

	return
}

func (t *Tasker) Remove(schedule task.Schedule) (err error) {
	_schedule := new(model.Schedule)
	_schedule.Id = schedule.Id()
	if _, de := t.schedule.Delete(_schedule); nil != de {
		err = de
	}

	return
}

func (t *Tasker) Running(id uint64, status task.Status, times uint32) (err error) {
	_task := new(model.Task)
	_task.Id = id
	_task.Status = status
	_task.Times = times
	if _, ue := t.task.Update(_task); nil != ue {
		err = ue
	}

	return
}

func (t *Tasker) Update(id uint64, status task.Status, runtime time.Time) (err error) {
	_task := new(model.Task)
	_task.Id = id
	_task.Status = status
	_task.Next = runtime
	if _, ue := t.task.Update(_task); nil != ue {
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
	if exists, ge := t.task.Get(_task); nil != ge {
		err = ge
	} else if !exists {
		err = exception.New().Message("任务不存在").Field(field.New("task", task)).Build()
	} else if _, ae := t.task.Archive(_task); nil != ae {
		err = ae
	}

	return
}

func (t *Tasker) Failed(_ task.Task) (err error) {
	return
}

func (t *Tasker) Run() (err error) {
	if tasks, re := t.task.GetsRunnable(); nil != re {
		err = re
	} else if 0 != len(*tasks) {
		all := *tasks
		t.runnable.Put(all[0], all[1:]...)
	}

	return
}
