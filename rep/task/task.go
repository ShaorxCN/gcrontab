package task

import (
	"gcrontab/constant"
	"gcrontab/entity/task"
	"gcrontab/model"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type TaskRep struct {
	db *gorm.DB
}

func NewTaskRep(db *gorm.DB) *TaskRep {
	if db == nil {
		db = model.DB()
	}

	return &TaskRep{db}
}

// FindActiveTasks 查找待运行的任务
func (r *TaskRep) FindActiveTasks(now time.Time) ([]*task.Task, error) {
	mts, err := model.FindActiveTasks(now)
	if err != nil {
		return nil, err
	}

	ts := make([]*task.Task, len(mts))
	var te *task.Task
	for i, t := range mts {
		if te, err = task.FromDBTaskModel(t); err != nil {
			logrus.Errorf("%v task status error:%v\n", t.ID, t.Status)
			return nil, err
		}
		ts[i] = te
	}

	return ts, nil
}

func (r *TaskRep) CreateTask(t *task.Task) error {
	m, err := t.ToDBTaskModel()
	if err != nil {
		return err
	}

	db := r.db
	if err := db.Create(m).Error; err != nil {
		logrus.Errorf("save task[%s] to db Failed:%v", m.Name, err)
		return err
	}
	return nil
}

func (r *TaskRep) FindTaskByID(id uuid.UUID) (*task.Task, error) {
	db := r.db
	dbTask := &model.DBTask{}
	err := db.Model(dbTask).Where("id = ? and status != ?", id, constant.STATUSDELDB).First(dbTask).Error
	if err != nil {
		return nil, err
	}

	return task.FromDBTaskModel(dbTask)
}
