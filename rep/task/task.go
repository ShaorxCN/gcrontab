package task

import (
	"fmt"
	"gcrontab/constant"
	"gcrontab/entity/task"
	"gcrontab/model"
	"gcrontab/rep/requestmodel"
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

	db := r.db
	var res []*model.DBTask
	err := db.Table(model.GetTaskTableName()).Where("next_runtime <= ? and Status = ?", now, constant.STATUSONDB).Find(&res).Error

	if err != nil {
		return nil, err
	}

	ts := make([]*task.Task, len(res))
	var te *task.Task
	for i, t := range res {
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

func (r *TaskRep) ModifyTaskTimeByID(id uuid.UUID, param *requestmodel.ModifyTask) error {
	db := r.db
	m := make(map[string]interface{})

	if !param.NextRuntimeUse.IsZero() {
		m["next_runtime"] = param.NextRuntimeUse
	}

	if !param.LastRuntimeUse.IsZero() {
		m["last_runtime"] = param.LastRuntimeUse
	}

	if !param.UpdateTimeUse.IsZero() {
		m["update_at"] = param.UpdateTimeUse
	}

	return db.Table(model.GetTaskTableName()).Where("id = ?", id).Updates(m).Error
}

func (r *TaskRep) FindTaskByName(name string) (*task.Task, error) {
	db := r.db
	t := new(model.DBTask)
	err := db.Model(t).Where("name = ? and status != ?", name, constant.STATUSDEL).First(t).Error
	if err != nil {
		return nil, err
	}
	return task.FromDBTaskModel(t)
}

// DeleteTaskByID 删除任务
func (r *TaskRep) DeleteTaskByID(id uuid.UUID) error {
	return r.db.Table(model.GetTaskTableName()).Where("id = ?", id).Update("status", constant.STATUSDEL).Error
}

func (r *TaskRep) FindTaskByNameWithOutStatus(name string) ([]uuid.UUID, error) {
	var taskIDs []uuid.UUID
	err := r.db.Table(model.GetTaskTableName()).Where("name ilike ?", fmt.Sprintf("%%%s%%", name)).Pluck("id", &taskIDs).Error
	return taskIDs, err
}

func (r *TaskRep) FindTasksByParam(p *requestmodel.TaskParams) ([]*task.Task, int, error) {
	limit := p.PageSize
	offset := (p.Page - 1) * p.PageSize

	db := r.db.Table(model.GetTaskTableName())

	if p.SortedBy != "" {
		if p.Order == constant.ASC || p.Order == constant.DESC {
			db = db.Order(fmt.Sprintf(" %s %s", p.SortedBy, p.Order))
		}
		db = db.Order(p.SortedBy)
	} else {
		db = db.Order("create_at DESC")
	}

	sqlStr, args := p.Task_buildQuery()
	db = db.Where(sqlStr, args...)
	var DBtasks []*model.DBTask
	var count int
	err := db.Count(&count).Error
	if err != nil {
		return nil, count, err
	}
	err = db.Limit(limit).Offset(offset).Find(&DBtasks).Error
	ts := make([]*task.Task, len(DBtasks))
	var te *task.Task
	for i, t := range DBtasks {
		if te, err = task.FromDBTaskModel(t); err != nil {
			logrus.Errorf("%v task status error:%v\n", t.ID, t.Status)
			return nil, count, err
		}
		ts[i] = te
	}
	return ts, count, err
}
