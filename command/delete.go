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
			index, err := cmd.Flags().GetIntSlice("index")
			if err != nil {
				cmd.Printf("Error getting flag 'index': %v\n", err)
				return
			}

			if len(index) == 0 {
				cmd.PrintErrln("flag --index is required")
				return
			}

			linkedTasksFlag, err := cmd.Flags().GetIntSlice("linkedTasks")
			if err != nil {
				cmd.Printf("Error getting flag 'linkedTasks': %v\n", err)
				return
			}

			// change linked tasks for task
			if len(linkedTasksFlag) != 0 {
				for _, i := range index {
					err = taskList.DeleteLinkedTasks(i, linkedTasksFlag...)
					if err != nil {
						cmd.Printf("Error changing linked tasks: %v\n", err)
						return
					}
				}
				return
			}

			err = taskList.Delete(index...)
			if err != nil {
				cmd.Printf("Error deleting task: %v\n", err)
				return
			}
		},
	}
	p := make([]int, 0)
	l := make([]int, 0)
	deleteCmd.Flags().IntSliceVarP(&p, "index", "i", p, "Task index")
	deleteCmd.Flags().IntSliceVarP(&l, "linkedTasks", "l", p, "Linked task index")

	return deleteCmd
}
