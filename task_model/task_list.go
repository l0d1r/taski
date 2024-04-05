package task_model

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"io/ioutil"
	"time"
)

const (
	ColorDefault = "\x1b[39m"

	ColorRed   = "\x1b[91m"
	ColorGreen = "\x1b[32m"
	ColorBlue  = "\x1b[94m"
	ColorGray  = "\x1b[90m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

func gray(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGray, s, ColorDefault)
}

type TaskList struct {
	store string
	list  []Task `json:"tasks"`
}

func NewTaskList(store string) *TaskList {
	return &TaskList{
		store: store,
	}
}

func (inst *TaskList) SetStore(store string) {
	inst.store = store
}

// Complete setup task by index in status done == true
func (inst *TaskList) Complete(index int) error {
	if index <= 0 || index > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	t := time.Now()

	inst.list[index-1].CompleteAt = &t
	inst.list[index-1].Done = true

	return inst.save()
}

// Delete method delete task from list by index
func (inst *TaskList) Delete(index int) error {
	if index <= 0 || index > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	if len(inst.list) == 1 {
		inst.list = nil
	} else {
		inst.list = append(inst.list[:index-1], inst.list[index:]...)
	}

	return inst.save()
}

// LoadFromStore loaded list of task from store
func (inst *TaskList) LoadFromStore() error {
	file, err := ioutil.ReadFile(inst.store)
	if err != nil {
		return err
	}

	if len(file) != 0 {
		err = json.Unmarshal(file, &inst.list)
		if err != nil {
			return err
		}
	}

	return nil
}

// Load method loading task descripton from file
func (inst *TaskList) Load(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, inst)
	if err != nil {
		return err
	}

	return nil
}

func (inst *TaskList) Add(name string, info string) error {
	inst.list = append(inst.list, *NewTask(name, info))
	return inst.save()
}
func (inst *TaskList) AddInfo(index int, info string) error {
	if index <= 0 || index > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}
	inst.list[index-1].Info = info

	return inst.save()
}

func (inst *TaskList) ViewInfo(index int) error {
	var (
		cells [][]*simpletable.Cell
	)
	if index <= 0 || index > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Info"},
		},
	}

	if inst.list[index-1].Info != "" {
		cells = append(
			cells,
			*&[]*simpletable.Cell{
				{Text: inst.list[index-1].Info},
			},
		)
	} else {
		cells = append(
			cells,
			*&[]*simpletable.Cell{
				{Text: ""},
			},
		)
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 1, Text: fmt.Sprintf("Task: %v", inst.list[index-1].Name)},
	}}

	table.Print()
	fmt.Print("\n")

	return nil
}
func (inst *TaskList) ViewTasks() error {
	var (
		cells [][]*simpletable.Cell
		tc    string
	)

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Status"},
			{Align: simpletable.AlignCenter, Text: "CreateAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	for idx, task := range inst.list {
		idx++
		t := blue(task.Name)
		s := red("no")
		if task.Done {
			t = green(task.Name)
			s = green("yes")
		}

		if task.CompleteAt != nil {
			tc = task.CompleteAt.Format(time.RFC822)
		}

		cells = append(
			cells,
			*&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", idx)},
				{Text: t},
				{Text: fmt.Sprintf("%s", s)},
				{Text: task.CreateAt.Format(time.RFC822)},
				{Text: tc},
			},
		)
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: "Tasks"},
	}}

	table.Print()
	fmt.Print("\n")

	return nil
}

func (inst *TaskList) save() error {
	b, err := json.Marshal(&inst.list)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(inst.store, b, 0644)
}
