package command

import (
	"github.com/spf13/cobra"
	"task/task_model"
)

func NewDeleteCmd(taskList *task_model.TaskList) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete task from list, usage only with flag -i (--index)",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			index, err := cmd.Flags().GetInt("index")
			if err != nil {
				cmd.Printf("Error getting flag 'index': %v\n", err)
				return
			}

			if index == 0 {
				cmd.PrintErrln("flag --index is required")
				return
			}

			err = taskList.Delete(index)
			if err != nil {
				cmd.Printf("Error deleting task: %v\n", err)
				return
			}
		},
	}

	deleteCmd.Flags().IntP("index", "i", 0, "Task index")

	return deleteCmd
}
