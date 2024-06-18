package service

import (
	"errors"
	"gcrontab/constant"
	"gcrontab/crontab"
	"gcrontab/entity/task"
	"gcrontab/interface/entity"
	"gcrontab/model"
	rep "gcrontab/rep/task"
	"gcrontab/utils"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type TaskService struct {
	ctx   *utils.ServiceContext
	db    *gorm.DB
	tasks []*task.Task
}

func dealNewTask(t *task.Task, operator, nickName string) {
	now := utils.Now()
	t.ID = entity.NewEntityKey(uuid.New().String(), task.TaskEntityType)
	t.CreateAt = now.Format(constant.TIMELAYOUT)
	t.Creater = operator
	t.CreaterName = nickName
	if t.Expired_time == 0 {
		t.Expired_time = model.TaskExpired
	}

	if t.Status == "" {
		t.Status = constant.STATUSOFF
	}
}

func NewTaskService(ctx *utils.ServiceContext, db *gorm.DB, tasks []*task.Task) *TaskService {
	if db == nil {
		db = model.DB()
	}
	if tasks == nil {
		return &TaskService{ctx, db, make([]*task.Task, 0, 10)}
	}
	return &TaskService{ctx, db, tasks}
}

func (ts *TaskService) CreateTask() error {
	if len(ts.tasks) == 0 {
		return errors.New("保存数据为空")
	}

	t := ts.tasks[0]

	dealNewTask(t, ts.ctx.Operator, ts.ctx.OperatorName)
	taskRep := rep.NewTaskRep(ts.db.New())
	taskRep.CreateTask(t)
	return nil
}

func (ts *TaskService) FindTaskByID(id uuid.UUID) (*task.Task, error) {
	taskRep := rep.NewTaskRep(ts.db.New())
	return taskRep.FindTaskByID(id)
}

func (ts *TaskService) RunTask(id uuid.UUID) error {
	taskRep := rep.NewTaskRep(ts.db.New())
	task, err := taskRep.FindTaskByID(id)
	if err != nil {
		return err
	}

	return crontab.ExecImmediately(task, ts.ctx.Operator)

}
