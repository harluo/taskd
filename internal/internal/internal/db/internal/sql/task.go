package sql

import (
	"fmt"
	"time"

	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/db/internal/get"
	"github.com/harluo/taskd/internal/internal/internal/internal/column"
	"github.com/harluo/taskd/internal/internal/internal/model"
	"github.com/harluo/xorm"
	"xorm.io/builder"
)

type Task struct {
	engine *xorm.Engine
	tx     *xorm.Tx

	table *model.Task
	tt    string
	cso   string
	cn    string
	csa   string
	ct    string
}

func NewTask(tx get.Tx) *Task {
	return &Task{
		engine: tx.Engine,
		tx:     tx.Tx,

		table: new(model.Task),
		tt:    tx.Engine.TableName(new(model.Task)),
		cso:   tx.Engine.ColumnName(column.Stop),
		cn:    tx.Engine.ColumnName(column.Next),
		csa:   tx.Engine.ColumnName(column.Status),
		ct:    tx.Engine.ColumnName(column.Times),
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
		t.cn: now, // 运行时间已到
	}.And(builder.Eq{
		t.csa: task.StatusCreated, // 刚创建的任务
	}.Or(builder.Eq{
		t.csa: task.StatusFailed, // 已经处于错误状态的任务
	}).Or(builder.Eq{
		t.csa: task.StatusStandby, // 已执行成功，但需继续执行的任务
	}))

	// 可被运行的条件二：任务已完成执行，但需要重启执行（可被循环执行的任务，达到下一次执行的条件）
	restarted := builder.Eq{ // 已经完成的任务，需要重新执行
		t.csa: task.StatusSuccess, // 已完成
	}.And(builder.Eq{
		t.ct: 0, // 次数被重置
	})

	// 可被运行的条件三：因各种问题中断执行
	interrupted := builder.NotNull{
		t.cso, // 结束时间有值
	}.And(builder.Lte{
		t.cso: now, // 超过最大运行时间段
	}.And(builder.Eq{
		t.csa: task.StatusRunning, // 运行中
	}.Or(builder.Eq{
		t.csa: task.StatusRetrying, // 重试中
	})))

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
	session.Join(scheduleTable, fmt.Sprintf("%s.id = %s.schedule", scheduleTable, taskTable))
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
		if dta, dte := session.Delete(deleted); nil != dte { // 删除计划本身
			err = dte
		} else if dsa, dse := t.deleteSchedule(session, task); nil != dse { // 删除对应的任务
			err = dse
		} else {
			affected = dta + dsa
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
