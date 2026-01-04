package table

import (
	"context"

	"github.com/harluo/migrate"
	"github.com/harluo/taskd/internal/internal/internal/model"
	"github.com/harluo/xorm"
)

type task struct {
	engine *xorm.Engine
}

func NewTask(engine *xorm.Engine) migrate.Migration {
	return &task{
		engine: engine,
	}
}

func (m *task) Upgrade(_ context.Context) error {
	return m.engine.CreateTables(new(model.Task))
}

func (m *task) Downgrade(_ context.Context) error {
	return m.engine.DropTables(new(model.Task))
}

func (*task) Id() uint64 {
	return 2025_05_07_11_22
}

func (*task) Description() string {
	return "创建任务调度细节表"
}

func (*task) Version() int {
	return 2
}
