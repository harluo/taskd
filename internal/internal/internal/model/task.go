package model

import (
	"fmt"
	"time"

	"github.com/goexl/gox"
	"github.com/goexl/task"
	"github.com/harluo/xorm"
)

type Task struct {
	xorm.Model `xorm:"extends"`

	// 计划
	Schedule uint64 `xorm:"BIGINT notnull index(target) default(0) comment(计划)" json:"schedule,omitempty"`
	// 开始执行时间
	Start time.Time `xorm:"DATETIME notnull index(next) default(CURRENT_TIMESTAMP) comment(开始时间)" json:"start,omitzero"` // nolint:lll
	// 下一次重试时间
	Next time.Time `xorm:"DATETIME notnull index(next) default(CURRENT_TIMESTAMP) comment(一下次执行时间)" json:"next,omitzero"` // nolint:lll
	// 结束时间
	Stop *time.Time `xorm:"DATETIME null index(next) comment(结束时间)" json:"stop,omitempty"`
	// 重试次数
	Times uint64 `xorm:"BIGINT notnull default(0) comment(重试次数)" json:"times,omitempty"`
	// 状态
	Status task.Status `xorm:"TINYINT notnull index(next) default(0) comment(状态，分别是：1、已创建；2、执行中；3、重试中；10、失败；20、成功)" json:"status,omitempty"` // nolint:lll
}

func (*Task) TableComment() string {
	return "任务调度细节"
}

func (t *Task) TaskId() (id string) {
	switch {
	case 0 != t.Id:
		id = gox.ToString(t.Id)
	default:
		id = fmt.Sprintf("%p", t)
	}

	return
}
