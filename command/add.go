package command

import (
	"github.com/spf13/cobra"
	"strings"
	"task/task_model"
	"time"
)

func NewAddCmd(taskList *task_model.TaskList) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a task",

		Run: func(cmd *cobra.Command, args []string) {
			info, err := cmd.Flags().GetString("description")
			if err != nil {
				cmd.Printf("Error getting flag 'description': %v\n", err)
				return
			}

			finishFlag, err := cmd.Flags().GetString("finish")
			if err != nil {
				cmd.Printf("Error getting flag 'finish': %v\n", err)
				return
			}

			finishFlagT, err := time.Parse(time.DateOnly, finishFlag)
			if err != nil {
				cmd.Printf("Error parsing flag 'finish': %v\n", err)
				return
			}

			err = taskList.Add(strings.Join(args, " "), info, &finishFlagT)
			if err != nil {
				cmd.Printf("Error adding task: %v\n", err)
				return
			}
		},
	}

	addCmd.Flags().StringP("description", "d", "", "Add additional info for task")
	addCmd.Flags().StringP("finish", "f", "", "Due finish date for task, format (2024-12-31), (yyyy-mm-dd)")
	return addCmd
}
