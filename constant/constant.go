package constant

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	TaskStatus = iota
	UserStatus
	DB2Enity
	Entity2DB
)

const (

	// header
	HEADEROPERATORID   = "OperatorID"
	HEADEROPERATORNAME = "Operator-Name"
	HEADEROPERATACCT   = "Operator-UserName"
	HEADEROPERATORROLE = "Operator-Role"
	HEADERTOKEN        = "Access-Token"

	// 时间格式

	TIMELAYOUT         = "2006-01-02 15:04:05"
	TIMELAYOUTWITHMILS = "2006-01-02 15:04:05.999"

	// task

	HTTPJOB        = "HTTP"
	EXECJON        = "COMMAND"
	STATUSPROCE    = "process"
	STATUSSUCC     = "success"
	STATUSFAIL     = "fail"
	HTTPMETHODGET  = "GET"
	HTTPMETHODPOST = "POST"
	POSTJSON       = "JSON"
	POSTFORM       = "FORM"

	POSTBODY    = "BODY"
	NOTIFYON    = "on"
	NOTIFYONDB  = 1
	NOTIFYOFF   = "off"
	NOTIFYOFFDB = 0
	STATUSON    = "on"
	STATUSONDB  = 1
	STATUSOFF   = "off"
	STATUSOFFDB = 0

	// query

	STARTTIME   = "startTime"
	ENDTIME     = "endTime"
	PAGE        = "page"
	PAGESIZE    = "pageSize"
	STATUS      = "status"
	NAME        = "name"
	CREATER     = "creater"
	SORTEDBY    = "sortedBy"
	ORDER       = "order"
	CREATERNAME = "createrName"
	LOGTASKID   = "taskID"
	TIMESTAMP   = "timeStamp"

	// role

	ANONYMOUS   = "anonymous"
	ADMIN       = "admin"
	ADMINDB     = 0
	TASKADMIN   = "taskAdmin"
	TASKADMINDB = 1
	USER        = "user"
	USERDB      = 3

	// DB

	// ASC 正序
	ASC = "ASC"
	// DESC 倒序
	DESC = "DESC"
	// 逻辑删除
	STATUSDEL      = "del"
	STATUSDELDB    = 2
	STATUSNORMAL   = "normal"
	STATUSNORMALDB = 0
)

var (
	Host string
)

func init() {
	host, err := os.Hostname()
	if err != nil {
		logrus.WithField("host", host).Errorf("get hostname failed:%v", err)
	}

	Host = host
}
