package service

import (
	tasklog "gcrontab/entity/task_log"
	rep "gcrontab/rep/task_log"
	"gcrontab/utils"
)

type TaskLogService struct {
	ctx *utils.ServiceContext
}

func NewTaskLogService(ctx *utils.ServiceContext) *TaskLogService {
	return &TaskLogService{ctx}
}

func (ts *TaskLogService) FindTaskLogByID() ([]*tasklog.TaskLog, int64, error) {
	return rep.FindTaskLogByID(ts.ctx.Query)
}
