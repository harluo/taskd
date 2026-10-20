package service

import (
	"context"

	"github.com/goexl/db"
	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/model"
)

func (t *Tasker) Add(required task.Schedule, optionals ...task.Schedule) (err error) {
	runtimes := make([]*model.Runtime, 0, 1+len(optionals))
	for _, _schedule := range append([]task.Schedule{required}, optionals...) {
		runtime := new(model.Runtime)
		runtime.Type = _schedule.Type()
		runtime.Next = _schedule.Next()
		runtime.Subtype = _schedule.Subtype()
		runtime.Target = _schedule.Target()
		runtime.Maximum = _schedule.Maximum()
		runtime.Timeout = db.NewDuration(_schedule.Timeout())
		runtime.Data = _schedule.Data()

		runtimes = append(runtimes, runtime)
	}
	if successes, ae := t.schedule.Add(context.Background(), runtimes[0], runtimes[1:]...); nil != ae {
		err = ae
	} else {
		t.runnable.Put((*successes)[0], (*successes)[1:]...)
	}

	return
}

func (t *Tasker) Remove(schedule task.Schedule) (err error) {
	_schedule := new(model.Schedule)
	_schedule.Id = schedule.Id()
	_schedule.Target = schedule.Target()
	_schedule.Type = schedule.Type()
	_schedule.Subtype = schedule.Subtype()
	_schedule.Maximum = schedule.Maximum()
	
	_schedule.Timeout = db.NewDuration(schedule.Timeout())
	_schedule.Data = schedule.Data()
	if _, de := t.schedule.Delete(context.Background(), _schedule); de != nil {
		err = de
	}

	return
}
