package main

import (
	"flag"
	"fmt"
	"os"
	"task/task_model"
)

func main() {
	var (
		taskList *task_model.TaskList
		err      error
	)
	argAdd := flag.String("add", "", "add new task in list, arguments awaits name of task")
	argInfo := flag.String("info", "", "task info")
	argAddInfo := flag.String("addInfo", "", "add note to task by index")
	argIndex := flag.Int("index", 0, "used for arg add info, required for addInfo")
	argGetInfo := flag.Int("getInfo", 0, "get task info by index")
	argComplete := flag.Int("complete", 0, "changed status of task on true, by id")
	argDelete := flag.Int("delete", 0, "delete task from list, by id")
	argStore := flag.String("store", "", "load task from file, argument await json file name ")
	flag.Parse()

	if argStore == nil || *argStore == "" {
		if s := os.Getenv("TASKI_STORE"); s == "" {
			store := fmt.Sprintf("%v/%v", os.Getenv("HOME"), ".taski/task.do.json")
			taskList = task_model.NewTaskList(store)
		}
	} else {
		taskList = task_model.NewTaskList(*argStore)
	}

	err = taskList.LoadFromStore()
	if err != nil {
		fmt.Printf("error load tasks from store: %v, err: %v", *argStore, err)
		os.Exit(1)
	}

	switch {
	case argAdd != nil && *argAdd != "":
		fmt.Printf("catch argument 'add': %v\n", *argAdd)
		if err = taskList.Add(*argAdd, *argInfo); err != nil {
			fmt.Printf("error adding task: %v", err)
			os.Exit(1)
		}
	case argAddInfo != nil && *argAddInfo != "":
		if argIndex == nil {
			fmt.Printf("missed argument index")
			os.Exit(1)
		}
		fmt.Printf("catch argument 'addInfo': %v\n", *argAddInfo)
		if err = taskList.AddInfo(*argIndex, *argAddInfo); err != nil {
			fmt.Printf("error adding task: %v", err)
			os.Exit(1)
		}
	case argGetInfo != nil && *argGetInfo != 0:
		if err = taskList.ViewInfo(*argGetInfo); err != nil {
			fmt.Printf("error get info: %v", err)
			os.Exit(1)
		}
	case argComplete != nil && *argComplete != 0:
		fmt.Printf("catch argument 'complete': %v\n", *argComplete)
		if err = taskList.Complete(*argComplete); err != nil {
			fmt.Printf("error complete task: %v", err)
			os.Exit(1)
		}
	case argDelete != nil && *argDelete != 0:
		fmt.Printf("catch argument 'delete': %v\n", *argDelete)
		if err = taskList.Delete(*argDelete); err != nil {
			fmt.Printf("error delete task: %v", err)
			os.Exit(1)
		}
	default:
		if err = taskList.ViewTasks(); err != nil {
			fmt.Printf("error view tasks: %v\n", err)
		}
	}
}
