package requestmodel

import (
	"time"

	"github.com/google/uuid"
)

// Params 是查询参数
type Params struct {
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
	CompanyCode string      `json:"companyCode,omitempty"`
	LogTaskID   string      `json:"taskID,omitempty"`
	TimeStamp   int64       `json:"timeStamp,omitempty"`
	TaskIDS     []uuid.UUID `json:"-"`
}
