package get

import (
	"github.com/goexl/log"
	"github.com/harluo/di"
	"github.com/harluo/schedule"
	"github.com/harluo/taskd/internal/internal/internal/core"
	"github.com/harluo/taskd/internal/internal/internal/db"
)

type Tasker struct {
	di.Get

	Schedule db.Schedule
	Task     db.Task

	Runnable  *core.Runnable
	Scheduler *schedule.Scheduler
	Logger    log.Logger
}
