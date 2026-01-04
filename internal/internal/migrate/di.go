package migrate

import (
	"github.com/goexl/gox"
	"github.com/harluo/di"
	"github.com/harluo/taskd/internal/internal/migrate/internal/table"
)

type Di = *gox.Empty

func init() {
	di.New().Instance().Put(
		table.NewSchedule,
		table.NewTask,
	).Group("migrations").Build().Apply()
}
