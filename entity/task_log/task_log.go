package tasklog

import (
	"gcrontab/constant"
	"gcrontab/interface/entity"
	"gcrontab/model"
	"gcrontab/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	TaskLogEntityType = "tbl_task_log"
)

type TaskLog struct {
	TimeStamp   int64  `json:"timeStamp,omitempty"`
	TaskName    string `json:"taskName,omitempty"`
	TaskID      string `json:"taskID,omitempty"`
	ResultCode  int    `json:"resultCode,omitempty"`
	Result      string `json:"result,omitempty"`
	Command     string `json:"command,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	TotalTime   int64  `json:"totalTime,omitempty"`
	Host        string `json:"host,omitempty"`
	Status      string `json:"status,omitempty"` // 执行状态  processing or  success or fail
	CompanyCode string `json:"companyCode,omitempty"`
	User        string `json:"user,omitempty"` // 如果是立即执行 这里会存储操作的用户名字
}

func (l *TaskLog) EntityKey() entity.Key {
	return entity.NewEntityKey(strconv.FormatInt(l.TimeStamp, 10), TaskLogEntityType)
}

// ToDBTaskLogModel 将tasklog实体转换为DBTaskLog数据库模型。
func (l *TaskLog) ToDBTaskLogModel() (*model.DBTaskLog, error) {

	start, err := time.ParseInLocation(constant.TIMELAYOUTWITHMILS, l.StartTime, utils.DefaultLocation)
	if err != nil {
		return nil, err
	}

	end, err := time.ParseInLocation(constant.TIMELAYOUTWITHMILS, l.EndTime, utils.DefaultLocation)
	if err != nil {
		return nil, err
	}

	d := &model.DBTaskLog{}

	taskID, err := uuid.Parse(l.TaskID)
	if err != nil {
		return nil, err
	}

	d.TimeStamp = l.TimeStamp
	d.TaskName = l.TaskName
	d.TaskID = taskID
	d.ResultCode = l.ResultCode
	d.Result = l.Result
	d.Command = l.Command
	d.StartTime = start
	d.EndTime = end

	d.TotalTime = l.TotalTime
	d.Host = l.Host

	d.Status = l.Status
	d.User = l.User

	return d, nil
}

// FromDBTaskLogModel 将任务日志Model转为实体。
func FromDBTaskLogModel(d *model.DBTaskLog) *TaskLog {
	t := new(TaskLog)
	t.TimeStamp = d.TimeStamp
	t.TaskName = d.TaskName
	t.TaskID = d.TaskID.String()
	t.ResultCode = d.ResultCode
	t.Result = d.Result
	t.Command = d.Command
	t.StartTime = d.StartTime.Format(constant.TIMELAYOUTWITHMILS)
	t.EndTime = d.EndTime.Format(constant.TIMELAYOUTWITHMILS)
	t.TotalTime = d.TotalTime
	t.Host = d.Host
	t.Status = d.Status
	t.User = d.User
	return t
}
