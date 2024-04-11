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
				cmd.Printf("Error getting flag 'description': %v\n", err)
				return
			}

			index, err := cmd.Flags().GetInt("index")
			if err != nil {
				cmd.Printf("Error getting flag 'index': %v\n", err)
				return
			}

			if index == 0 {
				cmd.PrintErrln("flag --index is required")
				return
			}

			if len(args) == 0 && description != "" {
				err = taskList.ChangeDescription(description, index)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
				return
			}

			statusFlag, err := cmd.Flags().GetBool("status")
			if err != nil {
				cmd.Printf("Error getting flag 'status': %v\n", err)
				return
			}

			if statusFlag {
				err = taskList.ChangeStatus(index)
				if err != nil {
					cmd.Printf("Error change status task: %v \n", err)
					return
				}
			}

			if len(args) != 0 {
				err = taskList.Change(args[0], description, index)
				if err != nil {
					cmd.Printf("Error change task: %v\n", err)
					return
				}
			}
		},
	}

	changeCmd.Flags().StringP("description", "d", "", "Add additional info for task")
	changeCmd.Flags().IntP("index", "i", 0, "Task index")
	changeCmd.Flags().BoolP("status", "s", false, "Change task status")

	return changeCmd
}
