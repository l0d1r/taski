package task_model

import (
	"crypto/sha1"
	"fmt"
	"time"
)

// Task contains data of task
type Task struct {
	Uuid          string              `json:"uuid"`
	LinkedTask    map[string]struct{} `json:"LinkedTask"`
	Name          string              `json:"Name"`
	Done          bool                `json:"Done"`
	CreateAt      time.Time           `json:"CreateAt"`
	CompleteAt    *time.Time          `json:"CompleteAt"`
	Info          string              `json:"Description"`
	DueFinishDate *time.Time          `json:"DueFinishDate"`
}

func NewTask(name string, info string, dueFinishDate *time.Time, linkedTask map[string]struct{}) *Task {
	uuid := sha1.New()
	uuid.Write([]byte(fmt.Sprintf("%s%s", name, time.Now().String())))

	return &Task{
		Uuid:          fmt.Sprintf("%x", uuid.Sum(nil)),
		LinkedTask:    linkedTask,
		Name:          name,
		Done:          false,
		CreateAt:      time.Now(),
		CompleteAt:    nil,
		Info:          info,
		DueFinishDate: dueFinishDate,
	}
}
