package main

import (
	"fmt"
	"os"
	"task/command"
	"task/task_model"
)

func CheckDefaultDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = nil
		err = os.Mkdir(fmt.Sprintf("%v/%v", os.Getenv("HOME"), ".taski"), 0775)
		if err != nil {
			return err
		}

		_, err = os.Create(fmt.Sprintf("%v/%v", os.Getenv("HOME"), ".taski/task.do.json"))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var taskList = new(task_model.TaskList)
	var lang = "RUS"
	var defaultPath = fmt.Sprintf("%v/%v", os.Getenv("HOME"), ".taski/task.do.json")

	if langEnv := os.Getenv("TASKI_LANG"); langEnv != "" {
		if langEnv != "ENG" && langEnv != "RUS" {
			fmt.Printf("only Russian (RUS) and English (ENG) language are available")
			os.Exit(1)
		}
		lang = langEnv
	}

	if path := os.Getenv("TASKI_STORE"); path != "" {
		defaultPath = path
	}

	if err := CheckDefaultDirExists(defaultPath); err != nil {
		fmt.Printf("Error check store directory: %v\n", err)
		os.Exit(1)
	}

	taskList = task_model.NewTaskList(
		defaultPath,
		lang,
	)

	if err := taskList.LoadFromStore(); err != nil {
		fmt.Printf("Error load file from store: %v \n", err)
		os.Exit(1)
	}

	rootCmd := command.NewRootCmd(taskList)
	rootCmd.AddCommand(command.NewAddCmd(taskList))
	rootCmd.AddCommand(command.NewViewCmd(taskList))
	rootCmd.AddCommand(command.NewDeleteCmd(taskList))
	rootCmd.AddCommand(command.NewChangeCmd(taskList))

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Error execute command: %v\n", err)
		os.Exit(1)
	}
}
