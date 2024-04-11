package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"task/task_model"
)

var viewStringRus = `Программа для управления, хранения задач. 

Каждая задача имеет статус,  время создания, время закрытия задачи, а также подробности.

- Время создания задачи устанавливается автоматически при создании задачи.

- Время закрытия задачи проставляется автоматически при изменении статуса задачи на "Выполнена", при смене статуса на "не выполнена" время выполнения убирается. 

- Дополнительная информация для задачи может быть пустым, а также задаваться разными способами, при создании задачи и также при помощи изменения.

В программе есть три, основных оператораю, для управления задачами в списке: add, delete, change и один для отображения списков.

Для подробной информации обратитесь к интересующей команде с флагом --help. (Пример: add --help)`

var viewStringEng = `Program for managing and storing tasks.

Each task has a status, creation time, task closing time, and details.

- The task creation time is set automatically when the task is created.

- The task closing time is entered automatically when the task status changes to “Completed”; when the status changes to “not completed”, the execution time is removed.

- Additional information for a task can be empty, and can also be specified in different ways, when creating a task and also by changing it.

The program has three main operators for managing tasks in a list: add, delete, change and one for displaying lists.

For detailed information, refer to the command of interest with the --help flag. (Example: add --help)`

func NewRootCmd(taskList *task_model.TaskList) *cobra.Command {
	return &cobra.Command{
		Use:   "taski",
		Short: "Application for contain list of tasks",
		Long:  viewStringEng,

		Run: func(cmd *cobra.Command, args []string) {
			if taskList.Language() == "RUS" {
				fmt.Println(viewStringRus)
				return
			}

			fmt.Println(viewStringEng)
		},
	}
}
