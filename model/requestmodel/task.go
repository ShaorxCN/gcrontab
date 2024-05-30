package requestmodel

import (
	"time"
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
