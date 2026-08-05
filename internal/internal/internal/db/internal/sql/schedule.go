package sql

import (
	"context"
	"time"

	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/db/internal/column"
	"github.com/harluo/taskd/internal/internal/internal/db/internal/get"
	"github.com/harluo/taskd/internal/internal/internal/model"
	"github.com/harluo/xorm"
	"xorm.io/builder"
)

type Schedule struct {
	engine *xorm.Engine
	tx     *xorm.Tx

	countBy *model.Schedule
	ts      string
	ci      string
	ct      string
}

func NewSchedule(gt get.Tx) *Schedule {
	return &Schedule{
		engine: gt.Engine,
		tx:     gt.Tx,

		countBy: new(model.Schedule),
		ts:      gt.Engine.TableName(new(model.Schedule)),
		ci:      gt.Engine.ColumnName(column.Id),
		ct:      gt.Engine.ColumnName(column.Target),
	}
}

func (s *Schedule) Add(
	ctx context.Context, runtime *model.Runtime, runtimes ...*model.Runtime,
) (successes *[]*model.Tasker, err error) {
	models := make([]*model.Tasker, 0, 1+len(runtimes))
	saves := append([]*model.Runtime{runtime}, runtimes...)
	if _, err = s.tx.Do(s.add(ctx, &saves, &models)); nil == err {
		successes = &models
	}

	return
}

func (s *Schedule) Get(ctx context.Context, schedule *model.Schedule, columns ...string) (bool, error) {
	return s.getSession(ctx, schedule, columns...).Get(schedule)
}

func (s *Schedule) Update(ctx context.Context, schedule *model.Schedule, columns ...string) (int64, error) {
	return s.engine.Context(ctx).ID(schedule.Id).MustCols(columns...).Update(schedule)
}

func (s *Schedule) Delete(ctx context.Context, schedule *model.Schedule) (int64, error) {
	return s.tx.Do(s.delete(ctx, schedule))
}

func (s *Schedule) delete(ctx context.Context, schedule *model.Schedule) func(session *xorm.Session) (int64, error) {
	return func(session *xorm.Session) (affected int64, err error) {
		session = session.Context(ctx)
		if ads, tasks, dse := s.deleteSchedule(session, schedule); nil != dse { // 删除计划本身
			err = dse
		} else if adt, dte := s.deleteTask(session, tasks); nil != dte { // 删除对应任务
			err = dte
		} else {
			affected = ads + adt
		}

		return
	}
}

func (s *Schedule) add(ctx context.Context, runtimes *[]*model.Runtime, successes *[]*model.Tasker) func(session *xorm.Session) (int64, error) {
	return func(session *xorm.Session) (affected int64, err error) {
		session = session.Context(ctx)
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
	tasks := make([]*model.Task, 0, len(*runtimes))
	for _, runtime := range *runtimes {
		_task := new(model.Task)
		_task.Schedule = runtime.Id
		_task.Next = runtime.Next
		_task.Status = task.StatusCreated

		now := time.Now()
		_task.Start = now
		if 0 != runtime.Timeout {
			stop := now.Add(runtime.Timeout)
			_task.Stop = &stop
		}

		tasks = append(tasks, _task)
	}

	if affected, err = session.Insert(tasks); nil == err {
		s.parseTasks(&tasks, runtimes, successes)
	}

	return
}

func (s *Schedule) deleteSchedule(
	session *xorm.Session, schedule *model.Schedule,
) (affected int64, tasks *[]*model.Task, err error) {
	schedules := make([]*model.Schedule, 0)
	if fe := session.Table(s.ts).Where(s.getsCond(schedule)).Find(&schedules); nil != fe {
		err = fe
	} else if affected, err = session.Delete(&schedules); nil == err {
		deletes := make([]*model.Task, 0, len(schedules))
		for _, matched := range schedules {
			deleted := new(model.Task)
			deleted.Schedule = matched.Id
			deletes = append(deletes, deleted)
		}
		tasks = &deletes
	}

	return
}

func (s *Schedule) deleteTask(session *xorm.Session, tasks *[]*model.Task) (affected int64, err error) {
	for _, task := range *tasks {
		var count int64
		if count, err = session.Delete(task); nil != err {
			affected += count
			return
		}
		affected += count
	}

	return
}

func (s *Schedule) getsCond(schedule *model.Schedule) (cond builder.Cond) {
	cond = builder.NewCond()
	if schedule.Id != 0 {
		cond = cond.And(builder.Eq{
			s.ci: schedule.Id,
		})
	}
	if schedule.Target != 0 {
		cond = cond.And(builder.Eq{
			s.ct: schedule.Target,
		})
	}

	return
}

func (s *Schedule) getSession(
	ctx context.Context,
	schedule *model.Schedule, columns ...string,
) *xorm.Session {
	return s.engine.Context(ctx).Table(s.ts).Cols(columns...).Where(s.getsCond(schedule))
}

func (s *Schedule) parseTasks(tasks *[]*model.Task, runtimes *[]*model.Runtime, successes *[]*model.Tasker) {
	for index, _task := range *tasks {
		success := new(model.Tasker)
		success.Id = _task.Id
		success.Start = _task.Start
		success.Next = _task.Next
		success.Stop = _task.Stop
		success.Times = _task.Times
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
