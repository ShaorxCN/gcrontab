package tasklog

import (
	"bytes"
	"fmt"
	"gcrontab/constant"
	tasklog "gcrontab/entity/task_log"
	tle "gcrontab/entity/task_log"
	"gcrontab/model"
	"gcrontab/rep/requestmodel"

	"github.com/google/uuid"
)

func SaveTaskLog(tl *tle.TaskLog) error {
	dbtl, err := tl.ToDBTaskLogModel()
	if err != nil {
		return err
	}

	return model.DB().Create(dbtl).Error
}

func UpdateTaskLog(tl *tle.TaskLog) error {
	dbtl, err := tl.ToDBTaskLogModel()
	if err != nil {
		return err
	}

	return model.DB().Save(dbtl).Error
}

func ModifyTaskTimeByID(id uuid.UUID, param *requestmodel.ModifyTask) error {

	return model.ModifyTaskTimeByID(id, param)
}

// FindTaskLogByPK 根据主键查找log
func FindTaskLogByPK(id uuid.UUID, timeStamp int64) (*tasklog.TaskLog, error) {
	tl := new(model.DBTaskLog)
	tl.TimeStamp = timeStamp
	tl.TaskID = id
	err := model.DB().Model(tl).First(tl).Error
	if err != nil {
		return nil, err
	}

	ret := tasklog.FromDBTaskLogModel(tl)
	return ret, err
}

func buildQuery4Log(p *requestmodel.Params) (string, []interface{}) {
	var buf bytes.Buffer

	args := make([]interface{}, 0, 3)
	buf.WriteString(" 1=1 ")

	if p.ID != uuid.Nil {
		buf.WriteString(" and task_id = ?")
		args = append(args, p.ID)
	}
	if !p.StartTime.IsZero() {
		buf.WriteString("   and start_time >= ?")
		args = append(args, p.StartTime)
	}

	if !p.EndTime.IsZero() {
		buf.WriteString(" and end_time <= ?")
		args = append(args, p.EndTime)
	}

	if p.Status != "" {
		buf.WriteString(" and status = ?")
		args = append(args, p.Status)
	}

	if p.TaskIDS != nil {
		buf.WriteString(" and task_id in (?)")
		args = append(args, p.TaskIDS)
	}

	return buf.String(), args
}

// FindTaskLogByID 根据任务id返回日志
func FindTaskLogByID(p *requestmodel.Params) ([]*tasklog.TaskLog, int64, error) {
	limit := p.PageSize
	offset := (p.Page - 1) * p.PageSize

	db := model.DB().Model(new(model.DBTaskLog))

	if p.SortedBy != "" {
		if p.Order == constant.ASC || p.Order == constant.DESC {
			db = db.Order(fmt.Sprintf(" %s %s", p.SortedBy, p.Order))
		} else {

			db = db.Order(p.SortedBy)
		}
	} else {
		db = db.Order("time_stamp DESC")
	}

	sqlStr, args := buildQuery4Log(p)
	db = db.Where(sqlStr, args...)

	var DBtaskLogs []*model.DBTaskLog
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, count, err
	}
	err = db.Limit(limit).Offset(offset).Find(&DBtaskLogs).Error

	return nil, count, err
}
