package task

import (
	"errors"
	"gcrontab/entity/task"
	rep "gcrontab/rep/task"
	"gcrontab/utils"
)

type TaskService struct {
	ctx   *utils.ServiceContext
	tasks []*task.Task
}

func NewTaskService(ctx *utils.ServiceContext, tasks []*task.Task) *TaskService {
	return &TaskService{ctx, tasks}
}

func (ts *TaskService) CreateTask() error {
	if len(ts.tasks) == 0 {
		return errors.New("保存数据为空")
	}

	rep.CreateTask(ts.tasks[0])
	return nil
}
