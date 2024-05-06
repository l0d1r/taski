package task_model

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
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
	w := strings.Fields(s)
	if len(s) > 50 {
		return fmt.Sprintf(
			"%s%s%s",
			ColorRed,
			fmt.Sprintf("%v%s\n%s%v%s", strings.Join(w[:len(w)/2], " "), ColorDefault, ColorRed, strings.Join(w[len(w)/2:], " "), ColorDefault),
			ColorDefault,
		)
	}
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	w := strings.Fields(s)
	if len(s) > 50 {
		return fmt.Sprintf(
			"%s%s%s",
			ColorGreen,
			fmt.Sprintf("%v%s\n%s%v%s", strings.Join(w[:len(w)/2], " "), ColorDefault, ColorGreen, strings.Join(w[len(w)/2:], " "), ColorDefault),
			ColorDefault,
		)
	}
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	w := strings.Fields(s)
	if len(s) > 50 {
		return fmt.Sprintf(
			"%s%s%s",
			ColorBlue,
			fmt.Sprintf("%v%s\n%s%v%s", strings.Join(w[:len(w)/2], " "), ColorDefault, ColorBlue, strings.Join(w[len(w)/2:], " "), ColorDefault),
			ColorDefault,
		)
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

func (inst *TaskList) Language() string {
	return inst.language
}

func (inst *TaskList) SetStore(store string) {
	inst.store = store
}

// ChangeStatus setup task by index in status done == true
func (inst *TaskList) ChangeStatus(index int) error {
	if index <= 0 || index > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	t := time.Now()

	if !inst.list[index-1].Done {
		inst.list[index-1].CompleteAt = &t
		inst.list[index-1].Done = true
	} else {
		inst.list[index-1].CompleteAt = nil
		inst.list[index-1].Done = false
	}

	return inst.save()
}

// Delete method delete task from list by index
func (inst *TaskList) Delete(index ...int) error {
	sort.Sort(sort.Reverse(sort.IntSlice(index)))
	for _, i := range index {
		if i <= 0 || i > len(inst.list) {
			return fmt.Errorf("invalid index\n")
		}

		if len(inst.list) == 1 {
			inst.list = nil
		} else {
			deleteTaskUuid := inst.list[i-1].Uuid

			for _, t := range inst.list {
				if _, ok := t.LinkedTask[deleteTaskUuid]; ok {
					delete(t.LinkedTask, deleteTaskUuid)
				}
			}

			inst.list = append(inst.list[:i-1], inst.list[i:]...)
		}
	}

	return inst.save()
}

func (inst *TaskList) DeleteLinkedTasks(taskIndex int, deletedTask ...int) error {
	if taskIndex <= 0 || taskIndex > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	for _, dt := range deletedTask {
		if dt <= 0 || dt > len(inst.list) {
			return fmt.Errorf("invalid index of linked task: %d\n", dt)
		}
		delete(inst.list[taskIndex-1].LinkedTask, inst.list[dt-1].Uuid)
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

// Load method loading task description from file
func (inst *TaskList) Load(filename string) error {
	var tasks []Task
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return err
	}

	inst.list = tasks

	return nil
}

func (inst *TaskList) Add(name string, info string, dueFinishDate *time.Time, linkedTask ...int) error {
	var linkedTasksUuid = make(map[string]struct{}, 0)

	for _, task := range linkedTask {
		if task <= 0 || task > len(inst.list) {
			return fmt.Errorf("invalid index\n")
		}
		if inst.list[task-1].Uuid == "" {
			uuid := sha1.New()
			uuid.Write([]byte(fmt.Sprintf("%s%s", name, time.Now().String())))
			inst.list[task-1].Uuid = fmt.Sprintf("%x", uuid.Sum(nil))
		}
		linkedTasksUuid[inst.list[task-1].Uuid] = struct{}{}
	}

	inst.list = append(inst.list, *NewTask(name, info, dueFinishDate, linkedTasksUuid))
	return inst.save()
}

func (inst *TaskList) AddLinkedTasks(task int, linkedTasks ...int) error {
	if task <= 0 || task > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	for _, tn := range linkedTasks {
		if tn <= 0 || tn > len(inst.list) {
			return fmt.Errorf("invalid index\n")
		}

		if inst.list[tn-1].Uuid == "" {
			uuid := sha1.New()
			uuid.Write([]byte(fmt.Sprintf("%s%s", inst.list[tn-1], time.Now().String())))
			inst.list[tn-1].Uuid = fmt.Sprintf("%x", uuid.Sum(nil))
		}

		inst.list[task-1].LinkedTask[inst.list[tn-1].Uuid] = struct{}{}
	}

	return inst.save()
}

func (inst *TaskList) Change(task, info string, idx int) error {
	if idx <= 0 || idx > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	inst.list[idx-1].Name = task
	inst.list[idx-1].Info = info

	return inst.save()
}

func (inst *TaskList) ChangeDescription(info string, idx int) error {
	if idx <= 0 || idx > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	inst.list[idx-1].Info = info

	return inst.save()
}

func (inst *TaskList) ChangeDueFinishDate(dft *time.Time, idx int) error {
	if idx <= 0 || idx > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}
	inst.list[idx-1].DueFinishDate = dft
	return inst.save()
}

func (inst *TaskList) ChangeLinkedTasks(numTask int, linkedTasks ...int) error {
	var linkedTaskMap = make(map[string]struct{}, 0)

	for _, task := range linkedTasks {
		if task <= 0 || task > len(inst.list) {
			return fmt.Errorf("invalid index\n")
		}

		if inst.list[task-1].Uuid == "" {
			uuid := sha1.New()
			uuid.Write([]byte(fmt.Sprintf("%s%s", inst.list[task-1], time.Now().String())))
			inst.list[task-1].Uuid = fmt.Sprintf("%x", uuid.Sum(nil))
		}

		linkedTaskMap[inst.list[task-1].Uuid] = struct{}{}
	}

	if numTask <= 0 || numTask > len(inst.list) {
		return fmt.Errorf("invalid index\n")
	}

	inst.list[numTask-1].LinkedTask = linkedTaskMap

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
				{Align: simpletable.AlignCenter, Text: "Description"},
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
		cells                [][]*simpletable.Cell
		taskTimeCompletedVal string
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
				{Align: simpletable.AlignCenter, Text: "Due Finish Date"},
				{Align: simpletable.AlignCenter, Text: "Linked Tasks"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 7, Text: "Tasks"},
		}}

	} else if inst.language == "RUS" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Номер"},
				{Align: simpletable.AlignCenter, Text: "Задача"},
				{Align: simpletable.AlignCenter, Text: "Статус"},
				{Align: simpletable.AlignCenter, Text: "Создано"},
				{Align: simpletable.AlignCenter, Text: "Дата/Время Выполнения"},
				{Align: simpletable.AlignCenter, Text: "Срок Окончания"},
				{Align: simpletable.AlignCenter, Text: "Связаные задачи"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 7, Text: "Задачи"},
		}}
	}

	for idx, task := range inst.list {
		var linkedTasksUuids []string
		var taskNameVal, taskStatusVal, taskFinishTimeVal string
		idx++
		if inst.language == "RUS" {
			taskNameVal = blue(task.Name)
			taskStatusVal = red("Не выполнена")
			if task.Done {
				taskNameVal = green(task.Name)
				taskStatusVal = green("Выполнена")
			}
		} else {
			taskNameVal = blue(task.Name)
			taskStatusVal = red("not completed")
			if task.Done {
				taskNameVal = green(task.Name)
				taskStatusVal = green("completed")
			}
		}

		if task.CompleteAt != nil {
			taskTimeCompletedVal = task.CompleteAt.Format(time.DateTime)
		} else {
			taskTimeCompletedVal = ""
		}

		if task.DueFinishDate != nil {
			taskFinishTimeVal = task.DueFinishDate.Format(time.DateOnly)
		} else {
			taskFinishTimeVal = ""
		}

		for taskUuid, _ := range task.LinkedTask {
			linkedTask, num := inst.taskByUuid(taskUuid)
			linkedTaskName := linkedTask.Name
			taskNameFields := strings.Fields(linkedTaskName)
			if len(taskNameFields) > 3 {
				linkedTaskName = fmt.Sprintf(
					"#%d: %s",
					*num+1,
					fmt.Sprintf("%v...", strings.Join(taskNameFields[:3], " ")),
				)
			} else {
				linkedTaskName = fmt.Sprintf("#%d: %s", *num+1, linkedTaskName)
			}

			if linkedTask.Done {
				linkedTaskName = green(linkedTaskName)
			} else {
				linkedTaskName = blue(linkedTaskName)
			}

			linkedTasksUuids = append(linkedTasksUuids, linkedTaskName)
		}

		cells = append(
			cells,
			*&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", idx)},
				{Text: taskNameVal},
				{Text: fmt.Sprintf("%s", taskStatusVal)},
				{Text: task.CreateAt.Format(time.DateTime)},
				{Text: taskTimeCompletedVal},
				{Text: taskFinishTimeVal},
				{Text: strings.Join(linkedTasksUuids, "\n ")},
			},
		)
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Print()
	fmt.Print("\n")

	return nil
}

func (inst *TaskList) ViewTask(index int) error {
	var (
		cells [][]*simpletable.Cell
		tc    string
	)

	task := inst.list[index-1]
	table := simpletable.New()

	if inst.language == "ENG" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "ID"},
				{Align: simpletable.AlignCenter, Text: "Name"},
				{Align: simpletable.AlignCenter, Text: "Status"},
				{Align: simpletable.AlignCenter, Text: "CreateAt"},
				{Align: simpletable.AlignCenter, Text: "CompletedAt"},
				{Align: simpletable.AlignCenter, Text: "Due Finish Date"},
				{Align: simpletable.AlignCenter, Text: "Linked Tasks"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 7, Text: "Tasks"},
		}}

	} else if inst.language == "RUS" {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Номер"},
				{Align: simpletable.AlignCenter, Text: "Задача"},
				{Align: simpletable.AlignCenter, Text: "Статус"},
				{Align: simpletable.AlignCenter, Text: "Создано"},
				{Align: simpletable.AlignCenter, Text: "Дата/Время Выполнения"},
				{Align: simpletable.AlignCenter, Text: "Срок Окончания"},
				{Align: simpletable.AlignCenter, Text: "Связаные задачи"},
			},
		}

		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 7, Text: "Задачи"},
		}}
	}

	var linkedTasks []string
	var t, s, ft string

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

	if task.DueFinishDate != nil {
		ft = task.DueFinishDate.Format(time.DateOnly)
	} else {
		ft = ""
	}

	for taskUuid, _ := range task.LinkedTask {
		linkedTask, num := inst.taskByUuid(taskUuid)
		linkedTaskName := linkedTask.Name
		tn := strings.Fields(linkedTaskName)
		if len(tn) > 3 {
			linkedTaskName = fmt.Sprintf(
				"#%d: %s",
				*num+1,
				fmt.Sprintf("%v...", strings.Join(tn[:3], " ")),
			)
		} else {
			linkedTaskName = fmt.Sprintf("#%d: %s", *num+1, linkedTaskName)
		}

		if linkedTask.Done {
			linkedTaskName = green(linkedTaskName)
		} else {
			linkedTaskName = blue(linkedTaskName)
		}

		linkedTasks = append(linkedTasks, linkedTaskName)
	}

	cells = append(
		cells,
		*&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index)},
			{Text: t},
			{Text: fmt.Sprintf("%s", s)},
			{Text: task.CreateAt.Format(time.DateTime)},
			{Text: tc},
			{Text: ft},
			{Text: strings.Join(linkedTasks, "\n ")},
		},
	)

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

func (inst *TaskList) taskByUuid(uuid string) (*Task, *int) {
	for num, task := range inst.list {
		if task.Uuid == uuid {
			return &task, &num
		}
	}
	return nil, nil
}
