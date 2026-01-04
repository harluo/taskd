package table

import (
	"context"

	"github.com/harluo/migrate"
	"github.com/harluo/taskd/internal/internal/internal/model"
	"github.com/harluo/xorm"
)

type schedule struct {
	engine *xorm.Engine
}

func NewSchedule(engine *xorm.Engine) migrate.Migration {
	return &schedule{
		engine: engine,
	}
}

func (s *schedule) Upgrade(_ context.Context) error {
	return s.engine.CreateTables(new(model.Schedule))
}

func (s *schedule) Downgrade(_ context.Context) error {
	return s.engine.DropTables(new(model.Schedule))
}

func (*schedule) Id() uint64 {
	return 2025_04_30_11_01
}

func (*schedule) Description() string {
	return "创建计划表"
}

func (*schedule) Version() int {
	return 2
}
