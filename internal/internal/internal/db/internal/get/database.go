package get

import (
	"github.com/goexl/log"
	"github.com/harluo/di"
	"github.com/harluo/xorm"
)

type Database struct {
	di.Get

	Logger log.Logger
	Engine *xorm.Engine
}
