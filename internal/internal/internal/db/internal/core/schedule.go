package core

import (
	"context"

	"github.com/harluo/taskd/internal/internal/internal/model"
)

type Schedule interface {
	Add(context.Context, *model.Runtime, ...*model.Runtime) (*[]*model.Tasker, error)

	Get(context.Context, *model.Schedule, ...string) (bool, error)

	Update(context.Context, *model.Schedule, ...string) (int64, error)

	Delete(context.Context, *model.Schedule) (int64, error)
}
