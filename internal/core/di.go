package core

import (
	"github.com/harluo/di"
	"github.com/harluo/taskd/internal/internal"
	"github.com/harluo/taskd/internal/internal/put"
)

func init() {
	di.New().Instance().Put(
		func(agent internal.Agent) put.Tasker {
			return put.Tasker{
				Tasker: agent.Tasker,
			}
		},
	).Build().Apply()
}
