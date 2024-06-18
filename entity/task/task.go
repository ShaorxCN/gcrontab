package task

import (
	"errors"
	"gcrontab/constant"
	"gcrontab/entity"
	entitier "gcrontab/interface/entity"
	"gcrontab/model"
	"gcrontab/utils"
	"time"

	"github.com/google/uuid"
)

const (
	// TaskEntityType entity type
	TaskEntityType = "task"
)

// Task 任务实体
type Task struct {
	entity.BaseEntity
	Name             string    `json:"name,omitempty" url:"name,omitempty"`
	IntervalDuration int       `json:"intervalDuration,omitempty" url:"intervalDuration,omitempty"` // n个周期的n
	Protocol         string    `json:"protocol,omitempty" url:"protocol,omitempty"`
	UnitOfInterval   string    `json:"unitOfInterval,omitempty" url:"unitOfInterval,omitempty"`
	Command          string    `json:"command,omitempty" url:"command,omitempty" `
	HTTPMethod       string    `json:"httpMethod,omitempty" url:"httpMethod,omitempty"`
	Param            string    `json:"param,omitempty" url:"param,omitempty"`
	PostType         string    `json:"postType,omitempty" url:"postType,omitempty"`
	Expired_time     int       `json:"expiredTime,omitempty" url:"expiredTime,omitempty"`
	Remark           string    `json:"remark,omitempty" url:"remark,omitempty"`
	NextRuntimeUse   time.Time `json:"-" url:"-"`
	LastRuntimeUse   time.Time `json:"-" url:"-"`
	Status           string    `json:"status,omitempty" url:"status,omitempty"`
	NextRuntime      string    `json:"nextRuntime,omitempty" url:"nextRuntime,omitempty"`
	LastRuntime      string    `json:"lastRuntime,omitempty" url:"lastRuntime,omitempty"`
	UpdateFlag       int8      `json:"-" url:"-"`
	UpdateTimeUse    time.Time `json:"-" url:"-"`
	UpdateID         string    `json:"updateID,omitempty" url:"updateID,omitempty"`
	CreaterName      string    `json:"createrName,omitempty" url:"createrName,omitempty"`
	Headers          string    `json:"headers,omitempty" url:"headers,omitempty"`
	Level            int       `json:"level,omitempty" url:"level,omitempty"`
}

// ToDBTaskModel 将task实体转换为DBTask数据库模型。
func (t *Task) ToDBTaskModel() (*model.DBTask, error) {
	d := &model.DBTask{}
	var err error
	var next, create, update, last time.Time
	next, err = time.ParseInLocation(constant.TIMELAYOUT, t.NextRuntime, utils.DefaultLocation)
	if err != nil {
		return nil, err
	}

	if t.LastRuntime != "" {
		last, err = time.ParseInLocation(constant.TIMELAYOUT, t.LastRuntime, utils.DefaultLocation)
		if err != nil {
			return nil, err
		}

		d.LastRuntime = &last
	}

	create, err = time.ParseInLocation(constant.TIMELAYOUT, t.CreateAt, utils.DefaultLocation)
	if err != nil {
		return nil, err
	}
	if t.UpdateAt != "" {
		update, err = time.ParseInLocation(constant.TIMELAYOUT, t.UpdateAt, utils.DefaultLocation)
		if err != nil {
			return nil, err
		}
		d.UpdateAt = &update
	}

	uid, err := uuid.Parse(t.ID.GetIDValue())
	if err != nil {
		return nil, err
	}
	d.ID = uid
	d.CreateAt = &create
	d.Creater = t.Creater
	d.Name = t.Name
	d.IntervalDuration = t.IntervalDuration
	d.UnitOfInterval = t.UnitOfInterval
	d.Protocol = t.Protocol
	d.Command = t.Command
	d.HTTPMethod = t.HTTPMethod
	d.Expired_time = t.Expired_time

	d.Remark = t.Remark

	switch t.Status {
	case constant.STATUSON:
		d.Status = constant.STATUSONDB
	case constant.STATUSOFF:
		d.Status = constant.STATUSOFFDB
	default:
		return nil, errors.New("status error")
	}

	d.NextRuntime = &next
	d.Level = t.Level
	d.UpdateID = t.UpdateID
	d.Param = t.Param
	d.PostType = t.PostType
	d.CreaterName = t.CreaterName
	d.Headers = t.Headers
	return d, nil
}

// FromDBTaskModel 将任务Model转为实体。
func FromDBTaskModel(d *model.DBTask) (*Task, error) {
	t := new(Task)
	t.ID = entitier.NewEntityKey(d.ID.String(), TaskEntityType)
	t.CreateAt = d.CreateAt.In(utils.DefaultLocation).Format(constant.TIMELAYOUT)
	if d.UpdateAt != nil {
		t.UpdateAt = d.UpdateAt.In(utils.DefaultLocation).Format(constant.TIMELAYOUT)
	}
	t.Creater = d.Creater
	t.Name = d.Name
	t.IntervalDuration = d.IntervalDuration
	t.UnitOfInterval = d.UnitOfInterval
	t.Protocol = d.Protocol
	t.Command = d.Command
	t.HTTPMethod = d.HTTPMethod

	t.Remark = d.Remark

	switch d.Status {
	case constant.STATUSONDB:
		t.Status = constant.STATUSON
	case constant.STATUSOFFDB:
		t.Status = constant.STATUSOFF
	default:
		return nil, errors.New("status error")
	}
	t.NextRuntime = d.NextRuntime.In(utils.DefaultLocation).Format(constant.TIMELAYOUT)
	t.NextRuntimeUse = *d.NextRuntime
	if d.LastRuntime != nil {
		t.LastRuntime = d.LastRuntime.In(utils.DefaultLocation).Format(constant.TIMELAYOUT)
	}
	t.Level = d.Level
	t.UpdateID = d.UpdateID
	t.Param = d.Param
	t.PostType = d.PostType
	t.CreaterName = d.CreaterName
	t.Headers = d.Headers
	return t, nil
}
