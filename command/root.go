package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"task/task_model"
)

func NewRootCmd(taskList *task_model.TaskList) *cobra.Command {
	return &cobra.Command{
		Use:   "taski",
		Short: "Application for contain list of tasks",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Application for contain list of tasks\n")
		},
	}
}
