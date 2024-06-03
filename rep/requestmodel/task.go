package requestmodel

import (
	"bytes"
	"fmt"
	"gcrontab/constant"
	"time"

	"github.com/google/uuid"
)

type ModifyTask struct {
	Name             string    `json:"name,omitempty" url:"name,omitempty"`
	IntervalDuration float64   `json:"intervalDuration,omitempty" url:"intervalDuration,omitempty"`
	UnitOfInterval   string    `json:"unitOfInterval,omitempty" url:"unitOfInterval,omitempty"`
	Protocol         string    `json:"protocol,omitempty" url:"protocol,omitempty"`
	Command          string    `json:"command,omitempty" url:"command,omitempty"`
	HTTPMethod       string    `json:"httpMethod,omitempty" url:"httpMethod,omitempty"`
	Param            string    `json:"param,omitempty" url:"param,omitempty"`
	PostType         string    `json:"postType,omitempty" url:"postType,omitempty"`
	TimeOut          int       `json:"timeOut,omitempty" url:"timeOut,omitempty"`
	Remark           string    `json:"remark,omitempty" url:"remark,omitempty"`
	NextRuntimeUse   time.Time `json:"-" url:"-"`
	LastRuntimeUse   time.Time `json:"-" url:"-"`
	Status           string    `json:"status,omitempty"  url:"status,omitempty"`
	UpdateFlag       int8      `json:"-" url:"-"`
	NextRuntime      string    `json:"nextRuntime,omitempty" url:"nextRuntime,omitempty"`
	UpdateTimeUse    time.Time `json:"-" url:"-"`
	Headers          string    `json:"headers,omitempty" url:"headers,omitempty"`
	Sign             string    `json:"sign,omitempty" url:"-" `
}

// TaskParams 可以参与查询条件的参数
type TaskParams struct {
	ID          uuid.UUID   `json:"-"`
	Page        int         `json:"page,omitempty"`
	PageSize    int         `json:"pageSize,omitempty"`
	Name        string      `json:"name,omitempty"`
	Status      string      `json:"status,omitempty"`
	Creater     string      `json:"creater,omitempty"`
	StartTime   time.Time   `json:"createTimeStart,omitempty"`
	EndTime     time.Time   `json:"createTimeEnd,omitempty"`
	SortedBy    string      `json:"sortedBy,omitempty"`
	Order       string      `json:"order,omitempty"`
	CreaterName string      `json:"createrName,omitempty"`
	LogTaskID   string      `json:"taskID,omitempty"`
	TimeStamp   int64       `json:"timeStamp,omitempty"`
	TaskIDS     []uuid.UUID `json:"-"`
}

func (p *TaskParams) Task_buildQuery() (string, []interface{}) {

	var buf bytes.Buffer

	args := make([]interface{}, 0, 3)
	if p.Status != "" {
		buf.WriteString("status = ?")
		args = append(args, p.Status)
	} else {
		buf.WriteString("status != ?")
		args = append(args, constant.STATUSDEL)
	}

	if p.Name != "" {
		buf.WriteString(" and name ilike ?")
		args = append(args, fmt.Sprintf("%%%s%%", p.Name))
	}

	if p.Creater != "" {
		buf.WriteString(" and creater = ?")
		args = append(args, p.Creater)
	}

	if p.Status != "" {
		buf.WriteString(" and status = ?")
		args = append(args, p.Status)
	}

	if p.CreaterName != "" {
		buf.WriteString(" and creater_name = ?")
		args = append(args, p.CreaterName)
	}

	if !p.StartTime.IsZero() && !p.EndTime.IsZero() {
		buf.WriteString(" and Create_at between ? and ?")
		args = append(args, p.StartTime, p.EndTime)
	} else {
		if !p.StartTime.IsZero() {
			buf.WriteString(" and Create_at >= ?")
			args = append(args, p.StartTime)
		}

		if !p.EndTime.IsZero() {
			buf.WriteString(" and Create_at <= ?")
			args = append(args, p.EndTime)
		}
	}

	return buf.String(), args
}
