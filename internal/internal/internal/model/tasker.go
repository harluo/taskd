package model

import (
	"time"

	"github.com/goexl/model"
	"github.com/goexl/task"
)

type Tasker struct {
	model.Base `xorm:"extends"`

	// 目标
	Target uint64 `json:"target,omitempty"`
	// 类型
	Type task.Type `json:"type,omitempty"`
	// 子类型
	Subtype task.Type `json:"subtype,omitempty"`
	// 超时时间
	Timeout time.Duration `json:"timeout,omitempty"`
	// 最大重试次数
	Maximum uint64 `json:"maximum,omitempty"`
	// 当前运行次数
	Times uint64 `json:"times,omitempty"`
	// 数据
	Data map[string]any `json:"data,omitempty"`

	// 开始执行时间
	Start time.Time `json:"start,omitzero"`
	// 下一次重试时间
	Next time.Time `json:"next,omitzero"`
	// 结束时间
	Stop *time.Time `json:"stop,omitempty"`
	// 状态
	Status task.Status `json:"status,omitempty"`
}
