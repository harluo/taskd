package model

import (
	"time"

	"github.com/goexl/task"
	"github.com/harluo/xorm"
)

type Schedule struct {
	xorm.Model `xorm:"extends"`

	// 目标
	Target uint64 `xorm:"BIGINT notnull index(target) default(0) comment(目标标识)" json:"target,omitempty"`
	// 类型
	Type task.Type `xorm:"TINYINT notnull index(next) default(0) comment(类型，分别是：1、表达式任务；2、固定时间任务；3、周期性任务；4、可被计算的任务；5、只执行一次的任务)" json:"type,omitempty"` // nolint:lll
	// 子类型
	Subtype task.Type `xorm:"SMALLINT notnull index(next) default(0) comment(子类型，根据应用自身识别)" json:"subtype,omitempty"` // nolint:lll
	// 超时时间
	Timeout time.Duration `xorm:"BIGINT notnull default(0) comment(任务超时时间)" json:"timeout,omitempty"`
	// 最大重试次数
	Maximum uint64 `xorm:"BIGINT notnull default(0) comment(最大重试次数)" json:"maximum,omitempty"`
	// 数据
	Data map[string]any `xorm:"LONGTEXT null comment(数据)" json:"data,omitempty"`
}

func (*Schedule) TableComment() string {
	return "计划"
}
