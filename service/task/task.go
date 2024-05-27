package task

import (
	"errors"
	"gcrontab/constant"
	"gcrontab/entity/task"
	"gcrontab/interface/entity"
	"gcrontab/model"
	rep "gcrontab/rep/task"
	"gcrontab/utils"

	"github.com/google/uuid"
)

type TaskService struct {
	ctx   *utils.ServiceContext
	tasks []*task.Task
}

func dealNewTask(t *task.Task, operator, nickName string) {
	now := utils.Now()
	t.ID = entity.NewEntityKey(uuid.New().String(), task.TaskEntityType)
	t.CreateAt = now.Format(constant.TIMELAYOUT)
	t.Creater = operator
	t.CreaterName = nickName
	if t.Expired_time == 0 {
		t.Expired_time = model.TimeOut
	}

	if t.Status == "" {
		t.Status = constant.STATUSON
	}
}

func NewTaskService(ctx *utils.ServiceContext, tasks []*task.Task) *TaskService {
	return &TaskService{ctx, tasks}
}

func (ts *TaskService) CreateTask() error {
	if len(ts.tasks) == 0 {
		return errors.New("保存数据为空")
	}

	t := ts.tasks[0]

	dealNewTask(t, ts.ctx.Operator, ts.ctx.OperatorName)

	rep.CreateTask(t)
	return nil
}
