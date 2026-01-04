package internal

import (
	"github.com/harluo/di"
	"github.com/harluo/taskd/internal/internal/migrate"
	"github.com/harluo/taskd/internal/internal/service"
)

type Agent struct {
	di.Get

	Tasker  *service.Tasker
	Migrate migrate.Di `optional:"true"`
}
