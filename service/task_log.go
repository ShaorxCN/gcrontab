package service

import (
	tasklog "gcrontab/entity/task_log"
	"gcrontab/model"
	rep "gcrontab/rep/task_log"
	"gcrontab/utils"

	"github.com/jinzhu/gorm"
)

type TaskLogService struct {
	ctx  *utils.ServiceContext
	db   *gorm.DB
	logs []*tasklog.TaskLog
}

func NewTaskLogService(ctx *utils.ServiceContext, db *gorm.DB, logs []*tasklog.TaskLog) *TaskLogService {
	if db == nil {
		db = model.DB()
	}
	if logs == nil {
		return &TaskLogService{ctx, db, make([]*tasklog.TaskLog, 0, 10)}
	}
	return &TaskLogService{ctx, db, logs}
}

func (ts *TaskLogService) FindTaskLogByID() ([]*tasklog.TaskLog, int64, error) {
	logRep := rep.NewTaskLogRep(ts.db.New())
	return logRep.FindTaskLogByID(ts.ctx.Query)
}
