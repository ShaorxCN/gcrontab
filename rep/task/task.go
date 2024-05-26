package task

import (
	"gcrontab/entity/task"
	"gcrontab/model"
	"time"

	"github.com/sirupsen/logrus"
)

// FindActiveTasks 查找待运行的任务
func FindActiveTasks(now time.Time) ([]*task.Task, error) {
	mts, err := model.FindActiveTasks(now)
	if err != nil {
		return nil, err
	}

	ts := make([]*task.Task, len(mts))

	for i, t := range mts {
		te := task.FromDBTaskModel(t)
		ts[i] = te
	}

	return ts, nil
}

func CreateTask(t *task.Task) error {
	m, err := t.ToDBTaskModel()
	if err != nil {
		return err
	}

	db := model.DB()
	if err := db.Create(m).Error; err != nil {
		logrus.Errorf("save task[%s] to db Failed:%v", m.Name, err)
		return err
	}
	return nil
}
