package kernel

import (
	"time"

	"github.com/goexl/task"
	"github.com/harluo/taskd/internal/internal/internal/constant"
	"github.com/harluo/taskd/internal/internal/internal/model"
)

var _ task.Task = (*Task)(nil)

type Task struct {
	task *model.Tasker
}

func NewTask(task *model.Tasker) *Task {
	return &Task{
		task: task,
	}
}

func (t *Task) Id() uint64 {
	return t.task.Id
}

func (t *Task) Target() uint64 {
	return t.task.Target
}

func (t *Task) Type() task.Type {
	return t.task.Type
}

func (t *Task) Subtype() task.Type {
	return t.task.Subtype
}

func (t *Task) Maximum() uint64 {
	return t.task.Maximum
}

func (t *Task) Times() uint64 {
	return t.task.Times
}

func (t *Task) Timeout() time.Duration {
	return t.task.Timeout
}

func (t *Task) Data() any {
	return t.task.Data[constant.FieldData]
}

func (t *Task) Next() time.Time {
	return t.task.Next
}
