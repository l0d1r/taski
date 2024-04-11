package command

import (
	"github.com/spf13/cobra"
	"task/task_model"
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

			err = taskList.Add(args[0], info)
			if err != nil {
				cmd.Printf("Error adding task: %v\n", err)
				return
			}
		},
	}

	addCmd.Flags().StringP("description", "d", "", "Add additional info for task")

	return addCmd
}
