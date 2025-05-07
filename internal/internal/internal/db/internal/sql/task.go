package sql

import (
	"fmt"
	"time"

	"github.com/goexl/id"
	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/db/internal/get"
	"github.com/harluo/taskd/internal/internal/internal/internal/column"
	"github.com/harluo/taskd/internal/internal/internal/model"
	"github.com/harluo/xorm"
	"xorm.io/builder"
)

type Task struct {
	id     id.Generator
	engine *xorm.Engine
	tx     *xorm.Tx
	table  *model.Task
}

func NewTask(tx get.Tx) *Task {
	return &Task{
		id:     tx.Id,
		engine: tx.DB,
		tx:     tx.Tx,
		table:  new(model.Task),
	}
}

func (t *Task) Add(task *model.Task) (affected int64, err error) {
	return t.engine.InsertOne(task)
}

func (t *Task) Get(task *model.Task, columns ...string) (bool, error) {
	return t.engine.Cols(columns...).Get(task)
}

func (t *Task) GetsRunnable(excludes ...*model.Task) (tasks *[]*model.Tasker, err error) {
	now := time.Now()

	// 可被运行的条件一：运行时间已到且状态是被认可可被重新执行
	timeout := builder.Lte{
		column.Next.String(): now, // 运行时间已到
	}.And(builder.Eq{
		column.Status.String(): task.StatusCreated, // 刚创建的任务
	}.Or(builder.Eq{
		column.Status.String(): task.StatusFailed, // 已经处于错误状态的任务
	}))

	// 可被运行的条件二：任务已完成执行，但需要重启执行（可被循环执行的任务，达到下一次执行的条件）
	restarted := builder.Eq{ // 已经完成的任务，需要重新执行
		column.Status.String(): task.StatusSuccess, // 已完成
	}.And(builder.Eq{
		column.Times.String(): 0, // 次数被重置
	})

	// 可被运行的条件三：因各种问题中断执行
	interrupted := builder.Lte{
		column.Stop.String(): now, // 超过最大运行时间段
	}.And(builder.Eq{
		column.Status.String(): task.StatusRunning, // 运行中
	}.Or(builder.Eq{
		column.Status.String(): task.StatusRetrying, // 重试中
	}))

	// 排除
	excludeTasks := builder.NewCond()
	for _, exclude := range excludes {
		excludeTasks = excludeTasks.And(builder.Neq{
			column.Id.String(): exclude.Id,
		})
	}

	entities := make([]*model.Tasker, 0)
	tasks = &entities
	condition := timeout.Or(restarted).Or(interrupted).And(excludeTasks)

	session := t.engine.Table(t.table).Where(condition)
	session.Limit(1024) // 最大取1024个数据

	taskTable := t.engine.TableName(t.table)
	scheduleTable := t.engine.TableName(new(model.Schedule))
	session.Join("INNER", scheduleTable, fmt.Sprintf("%s.id = %s.schedule", scheduleTable, taskTable))
	err = session.Find(tasks)

	return
}

func (t *Task) Update(task *model.Task, columns ...string) (int64, error) {
	return t.engine.ID(task.Id).MustCols(columns...).Update(task)
}

func (t *Task) Archive(task *model.Task) (int64, error) {
	return t.tx.Do(t.delete(task))
}

func (t *Task) Delete(task *model.Task) (int64, error) {
	return t.tx.Do(t.delete(task))
}

func (t *Task) delete(task *model.Task) func(session *xorm.Session) (int64, error) {
	return func(session *xorm.Session) (affected int64, err error) {
		deleted := new(model.Task)
		deleted.Id = task.Id
		if adt, dse := session.Delete(deleted); nil != dse { // 删除计划本身
			err = dse
		} else if ads, dte := t.deleteSchedule(session, task); nil != dte { // 删除对应的任务
			err = dte
		} else {
			affected = adt + ads
		}

		return
	}
}

func (t *Task) deleteSchedule(session *xorm.Session, task *model.Task) (affected int64, err error) {
	deleted := new(model.Schedule)
	deleted.Id = task.Schedule
	affected, err = session.Delete(deleted)

	return
}
