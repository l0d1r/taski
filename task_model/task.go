package task_model

import "time"

type Task struct {
	Name       string     `json:"Name"`
	Done       bool       `json:"Done"`
	CreateAt   time.Time  `json:"CreateAt"`
	CompleteAt *time.Time `json:"CompleteAt"`
	Info       string     `json:"Description"`
}

func NewTask(name string, info string) *Task {
	return &Task{
		Name:       name,
		Done:       false,
		CreateAt:   time.Now(),
		CompleteAt: nil,
		Info:       info,
	}
}
