package task_model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alexeyco/simpletable"
)

const (
	ColorDefault = "\x1b[39m"

	ColorRed   = "\x1b[91m"
	ColorGreen = "\x1b[32m"
	ColorBlue  = "\x1b[94m"
	ColorGray  = "\x1b[90m"
)

func red(s string) string {
	if len(s) > 30 {
		return fmt.Sprintf("%s%s\n%s%s%s", ColorRed, s[:int(len(s)/2)], ColorRed, s[int(len(s)/2):], ColorDefault)
	}
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	if len(s) > 30 {
		return fmt.Sprintf("%s%s\n %s%s%s", ColorGreen, s[:int(len(s)/2)], ColorGreen, s[int(len(s)/2):], ColorDefault)
	}
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	if len(s) > 30 {
		return fmt.Sprintf("%s%s\n%s,%s%s", ColorBlue, s[:int(len(s)/2)], ColorBlue, s[int(len(s)/2):], ColorDefault)
	}
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

func gray(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGray, s, ColorDefault)
}

type TaskList struct {
	store    string
	language string
	list     []Task `json:"tasks"`
}

func NewTaskList(store string, language string) *TaskList {
	return &TaskList{
		store:    store,
		language: language,
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

	if inst.language == "ENG" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Info"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 1, Text: fmt.Sprintf("Task: %v", green(inst.list[index-1].Name))},
		}}

	} else if inst.language == "RUS" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Подробности"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 1, Text: fmt.Sprintf("Задача: %v", green(inst.list[index-1].Name))},
		}}
	}

	if inst.list[index-1].Info != "" {
		cells = append(
			cells,
			*&[]*simpletable.Cell{
				{Text: green(inst.list[index-1].Info)},
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

	if inst.language == "ENG" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "ID"},
				{Align: simpletable.AlignCenter, Text: "Name"},
				{Align: simpletable.AlignCenter, Text: "Status"},
				{Align: simpletable.AlignCenter, Text: "CreateAt"},
				{Align: simpletable.AlignCenter, Text: "CompletedAt"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: "Tasks"},
		}}

	} else if inst.language == "RUS" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Номер"},
				{Align: simpletable.AlignCenter, Text: "Задача"},
				{Align: simpletable.AlignCenter, Text: "Статус"},
				{Align: simpletable.AlignCenter, Text: "Создано"},
				{Align: simpletable.AlignCenter, Text: "Выполненно"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: "Задачи"},
		}}
	}

	for idx, task := range inst.list {
		var t, s string
		idx++
		if inst.language == "RUS" {
			t = blue(task.Name)
			s = red("невыполненна")
			if task.Done {
				t = green(task.Name)
				s = green("выполненна")
			}
		} else {
			t = blue(task.Name)
			s = red("no")
			if task.Done {
				t = green(task.Name)
				s = green("yes")
			}
		}

		if task.CompleteAt != nil {
			tc = task.CompleteAt.Format(time.RFC822)
		} else {
			tc = ""
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
