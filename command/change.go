package command

import (
	"github.com/spf13/cobra"
	"task/task_model"
)

func NewChangeCmd(taskList *task_model.TaskList) *cobra.Command {
	changeCmd := &cobra.Command{
		Use:   "change",
		Short: "Change task name or info, usage only with flag -i (--index) ",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			description, err := cmd.Flags().GetString("description")
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			index, err := cmd.Flags().GetInt("index")
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if index == 0 {
				cmd.PrintErrln("flag --index is required")
				return
			}

			if len(args) == 0 || args[0] == "" {
				err = taskList.ChangeDescription(description, index)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
				return
			}

			err = taskList.Change(args[0], description, index)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	changeCmd.Flags().StringP("description", "d", "", "Add additional info for task")
	changeCmd.Flags().IntP("index", "i", 0, "Task index")

	return changeCmd
}
