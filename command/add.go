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
			index, err := cmd.Flags().GetInt("index")
			if err != nil {
				cmd.Printf("Error getting flag 'index': %v\n", err)
				return
			}

			linkedTasksFlag, err := cmd.Flags().GetIntSlice("linkedTasks")
			if err != nil {
				cmd.Printf("Error getting flag 'linkedTasks': %v\n", err)
				return
			}

			// add linked tasks for task
			if len(linkedTasksFlag) != 0 && index != 0 {
				err = taskList.AddLinkedTasks(index, linkedTasksFlag...)
				if err != nil {
					cmd.Printf("Error adding linked tasks: %v\n", err)
					return
				}
				return
			}

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

			finishFlagT := new(time.Time)
			if finishFlag != "" {
				*finishFlagT, err = time.Parse(time.DateOnly, finishFlag)
				if err != nil {
					cmd.Printf("Error parsing flag 'finish': %v\n", err)
					return
				}
			} else {
				finishFlagT = nil
			}

			err = taskList.Add(strings.Join(args, " "), info, finishFlagT, linkedTasksFlag...)
			if err != nil {
				cmd.Printf("Error adding task: %v\n", err)
				return
			}

		},
	}

	p := make([]int, 0)

	addCmd.Flags().StringP("description", "d", "", "Add additional info for task")
	addCmd.Flags().IntSliceVarP(&p, "linkedTasks", "l", p, "Linked task index")
	addCmd.Flags().IntP("index", "i", 0, "Task index")
	addCmd.Flags().StringP("finish", "f", "", "Due finish date for task, format (2024-12-31), (yyyy-mm-dd)")
	return addCmd
}
