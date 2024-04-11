package command

import (
	"github.com/spf13/cobra"
	"task/task_model"
)

func NewViewCmd(taskList *task_model.TaskList) *cobra.Command {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View task list",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			indexFlag, err := cmd.Flags().GetInt("index")
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			descriptionFlag, err := cmd.Flags().GetBool("description")
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if descriptionFlag {
				if indexFlag == 0 {
					cmd.PrintErrln("You must specify the index")
					return
				}

				err = taskList.ViewInfo(indexFlag)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				return
			}

			if indexFlag != 0 {
				err = taskList.ViewTask(indexFlag)
				if err != nil {
					cmd.PrintErrln(err)
				}
				return
			}

			err = taskList.ViewTasks()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	viewCmd.Flags().IntP("index", "i", 0, "task index")
	viewCmd.Flags().BoolP("description", "d", false, "task description")

	return viewCmd
}
