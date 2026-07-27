package core

import (
	"github.com/harluo/taskd/internal/internal/internal/model"
)

type Schedule interface {
	Add(*model.Runtime, ...*model.Runtime) (*[]*model.Tasker, error)

	Get(*model.Schedule, ...string) (bool, error)

	Update(*model.Schedule, ...string) (int64, error)

	Delete(*model.Schedule) (int64, error)
}
