package task

import (
	"gcrontab/entity/task"
	"gcrontab/model"
	"time"
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
