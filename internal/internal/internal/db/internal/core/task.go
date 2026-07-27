package core

import (
	"github.com/harluo/taskd/internal/internal/internal/model"
)

type Task interface {
	Add(task *model.Task) (int64, error)

	Get(task *model.Task, columns ...string) (bool, error)

	GetsRunnable(excludes ...*model.Task) (*[]*model.Tasker, error)

	Update(task *model.Task, columns ...string) (int64, error)

	Archive(task *model.Task) (int64, error)

	Delete(task *model.Task) (int64, error)
}
