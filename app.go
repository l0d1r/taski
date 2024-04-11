package main

import (
	"fmt"
	"os"
	"task/command"
	"task/task_model"
)

func main() {
	var taskList = new(task_model.TaskList)
	var lang = "RUS"

	if langEnv := os.Getenv("TASKI_LANG"); langEnv != "" {
		if langEnv != "ENG" && langEnv != "RUS" {
			fmt.Printf("only Russian (RUS) and English (ENG) language are available")
			os.Exit(1)
		}

		lang = langEnv
	}

	if os.Getenv("TASKI_STORE") == "" {
		taskList = task_model.NewTaskList(
			fmt.Sprintf("%v/%v", os.Getenv("HOME"), ".taski/task.do.json"),
			lang,
		)
	} else {
		taskList = task_model.NewTaskList(
			os.Getenv("TASKI_STORE"),
			lang,
		)
	}

	if err := taskList.LoadFromStore(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd := command.NewRootCmd(taskList)
	rootCmd.AddCommand(command.NewAddCmd(taskList))
	rootCmd.AddCommand(command.NewViewCmd(taskList))
	rootCmd.AddCommand(command.NewDeleteCmd(taskList))
	rootCmd.AddCommand(command.NewChangeCmd(taskList))

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
