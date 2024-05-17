package tasklog

import (
	tle "gcrontab/entity/task_log"
	"gcrontab/model"
)

func SaveTaskLog(tl *tle.TaskLog) error {
	dbtl, err := tl.ToDBTaskLogModel()
	if err != nil {
		return err
	}

	return model.InsertTaskLog(dbtl)
}

func UpdateTaskLog(tl *tle.TaskLog) error {
	dbtl, err := tl.ToDBTaskLogModel()
	if err != nil {
		return err
	}

	return model.UpdateTaskLog(dbtl)
}
