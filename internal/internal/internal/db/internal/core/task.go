package core

import (
	"context"

	"github.com/harluo/taskd/internal/internal/internal/model"
)

type Task interface {
	Add(context.Context,  *model.Task) (int64, error)

	Get(context.Context,  *model.Task,  ...string) (bool, error)

	GetsRunnable(context.Context,  ...*model.Task) (*[]*model.Tasker, error)

	Update(context.Context,  *model.Task,  ...string) (int64, error)

	Archive(context.Context,  *model.Task) (int64, error)

	Delete(context.Context,  *model.Task) (int64, error)
}
