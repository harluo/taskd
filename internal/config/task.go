package config

import (
	"time"

	"github.com/harluo/config"
)

// Task 任务配置
type Task struct {
	// 任务最长执行时间
	Lifetime time.Duration `default:"30m" json:"lifetime,omitempty"`
}

func newTask(getter config.Getter) (task *Task, err error) {
	task = new(Task)
	err = getter.Get(&struct {
		Task *Task `json:"task,omitempty"`
	}{
		Task: task,
	})

	return
}
