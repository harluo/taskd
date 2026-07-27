package sql

import (
	"time"

	"github.com/goexl/id"
	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/db/internal/get"
	"github.com/harluo/taskd/internal/internal/internal/model"
	"github.com/harluo/xorm"
)

type Schedule struct {
	id     id.Generator
	engine *xorm.Engine
	tx     *xorm.Tx
}

func NewSchedule(database get.Tx) *Schedule {
	return &Schedule{
		id:     database.Id,
		engine: database.DB,
		tx:     database.Tx,
	}
}

func (s *Schedule) Add(runtime *model.Runtime, runtimes ...*model.Runtime) (successes *[]*model.Tasker, err error) {
	models := make([]*model.Tasker, 0, 1+len(runtimes))
	saves := append([]*model.Runtime{runtime}, runtimes...)
	if _, err = s.tx.Do(s.add(&saves, &models)); nil == err {
		successes = &models
	}

	return
}

func (s *Schedule) Get(schedule *model.Schedule, columns ...string) (bool, error) {
	return s.engine.Cols(columns...).Get(schedule)
}

func (s *Schedule) Update(schedule *model.Schedule, columns ...string) (int64, error) {
	return s.engine.ID(schedule.Id).MustCols(columns...).Update(schedule)
}

func (s *Schedule) Delete(schedule *model.Schedule) (int64, error) {
	return s.tx.Do(s.delete(schedule))
}

func (s *Schedule) delete(schedule *model.Schedule) func(session *xorm.Session) (int64, error) {
	return func(session *xorm.Session) (affected int64, err error) {
		deleted := new(model.Schedule)
		deleted.Id = schedule.Id
		if ads, dse := session.Delete(deleted); nil != dse { // 删除计划本身
			err = dse
		} else if adt, dte := s.deleteTask(session, schedule); nil != dte { // 删除对应的任务
			err = dte
		} else {
			affected = ads + adt
		}

		return
	}
}

func (s *Schedule) add(runtimes *[]*model.Runtime, successes *[]*model.Tasker) func(session *xorm.Session) (int64, error) {
	return func(session *xorm.Session) (affected int64, err error) {
		schedules := make([]any, 0, len(*runtimes))
		for _, runtime := range *runtimes {
			schedule := &runtime.Schedule
			schedules = append(schedules, schedule)
		}

		if ais, ie := session.Insert(schedules); nil != ie {
			err = ie
		} else if aat, ate := s.addTasks(session, runtimes, successes); nil != ate {
			err = ate
		} else {
			affected = ais + aat
		}

		return
	}
}

func (s *Schedule) addTasks(
	session *xorm.Session,
	runtimes *[]*model.Runtime, successes *[]*model.Tasker,
) (affected int64, err error) {
	tasks := make([]any, 0, len(*runtimes))
	for _, runtime := range *runtimes {
		_task := new(model.Task) // !不用设置标识，通过事件注入
		_task.Schedule = runtime.Id
		_task.Next = runtime.Next
		_task.Status = task.StatusCreated

		now := time.Now()
		_task.Start = now
		_task.Stop = now.Add(runtime.Timeout)

		tasks = append(tasks, _task)
	}

	if affected, err = session.Insert(tasks); nil == err {
		s.parseTasks(&tasks, runtimes, successes)
	}

	return
}

func (s *Schedule) deleteTask(session *xorm.Session, schedule *model.Schedule) (affected int64, err error) {
	deleted := new(model.Task)
	deleted.Schedule = schedule.Id
	affected, err = session.Delete(deleted)

	return
}

func (s *Schedule) parseTasks(tasks *[]any, runtimes *[]*model.Runtime, successes *[]*model.Tasker) {
	for index, _task := range *tasks {
		if converted, ok := _task.(*model.Task); ok {
			success := new(model.Tasker)
			success.Id = converted.Id
			success.Start = converted.Start
			success.Next = converted.Next
			success.Stop = converted.Stop
			success.Times = converted.Times
			success.Status = task.StatusCreated

			schedule := (*runtimes)[index]
			success.Target = schedule.Target
			success.Type = schedule.Type
			success.Subtype = schedule.Subtype
			success.Maximum = schedule.Maximum
			success.Timeout = schedule.Timeout
			success.Data = schedule.Data

			*successes = append(*successes, success)
		}
	}
}
