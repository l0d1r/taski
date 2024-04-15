package command

import (
	"github.com/spf13/cobra"
	"strings"
	"task/task_model"
	"time"
)

func NewChangeCmd(taskList *task_model.TaskList) *cobra.Command {
	changeCmd := &cobra.Command{
		Use:   "change",
		Short: "Change task name or info, usage only with flag -i (--index) ",
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

			finishFlag, err := cmd.Flags().GetString("finish")
			if err != nil {
				cmd.Printf("Error getting flag 'finish': %v\n", err)
				return
			}

			description, err := cmd.Flags().GetString("description")
			if err != nil {
				cmd.Printf("Error getting flag 'description': %v\n", err)
				return
			}

			statusFlag, err := cmd.Flags().GetBool("status")
			if err != nil {
				cmd.Printf("Error getting flag 'status': %v\n", err)
				return
			}

			if finishFlag != "" {
				finishFlagT, err := time.Parse(time.DateOnly, finishFlag)
				if err != nil {
					cmd.Printf("Error parsing flag 'finish': %v\n", err)
					return
				}

				err = taskList.ChangeDueFinishDate(&finishFlagT, index)
				if err != nil {
					cmd.Printf("Error changing due date: %v\n", err)
					return
				}
			}

			if statusFlag {
				err = taskList.ChangeStatus(index)
				if err != nil {
					cmd.Printf("Error change status task: %v \n", err)
					return
				}
			}

			if len(args) == 0 && description != "" {
				err = taskList.ChangeDescription(description, index)
				if err != nil {
					cmd.Printf("Error change description: %v\n", err)
					return
				}
				return
			}

			if len(args) != 0 {
				err = taskList.Change(strings.Join(args, " "), description, index)
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
	changeCmd.Flags().StringP("finish", "f", "", "Due finish date for task, format (2024-12-31), (yyyy-mm-dd)")

	return changeCmd
}
