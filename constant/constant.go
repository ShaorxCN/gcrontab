package constant

const (

	// header

	HEADEROPERATORNAME = "Operator-Name"
	HEADEROPERATOR     = "Operator"
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

	POSTBODY  = "BODY"
	NOTIFYON  = "on"
	NOTIFYOFF = "off"
	STATUSON  = "on"
	STATUSOFF = "off"

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
	COMPANYCODE = "companyCode"
	LOGTASKID   = "taskID"
	TIMESTAMP   = "timeStamp"

	// role

	ANONYMOUS = "anonymous"
	ADMIN     = "admin"
	USER      = "user"
	TASKADMIN = "taskAdmin"

	// DB

	// ASC 正序
	ASC = "ASC"
	// DESC 倒序
	DESC = "DESC"
	// 逻辑删除
	STATUSDEL = "del"
)