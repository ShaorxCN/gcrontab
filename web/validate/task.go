package validate

import (
	"encoding/json"
	"fmt"
	"gcrontab/constant"
	"gcrontab/custom"
	"gcrontab/entity/task"
	"gcrontab/utils"
	"net/http"
	"strings"
	"time"
)

var (
	units     = []string{"second", "minute", "hour", "day", "month", "week"}
	protocols = []string{constant.HTTPJOB}
	postSlice = []string{"JSON", "BODY", ""}
	methods   = []string{"GET", "POST"}
)

// CheckCreateTaskRequest 创建任务request的基础校验
func CheckCreateTaskRequest(in *task.Task) error {
	switch {
	case in.Name == "" || len(in.Name) > 512:
		return custom.ParamErrorReturn(name)
	case in.IntervalDuration == 0:
		return custom.ParamErrorReturn(intervalDuration)
	case in.Command == "":
		return custom.ParamErrorReturn(command)
	}

	if !utils.StrInSlice(in.UnitOfInterval, units) {
		return custom.ParamErrorReturn(unitOfInterval)
	}

	if !utils.StrInSlice(in.Protocol, protocols) {
		return custom.ParamErrorReturn(protocol)

	}

	if in.Param != "" && len(in.Param) > 4096 {
		return custom.ParamErrorReturn(param)
	}

	if in.HTTPMethod != "" && !utils.StrInSlice(in.HTTPMethod, methods) {
		return custom.ParamErrorReturn(httpMethod)

	}

	if !utils.StrInSlice(in.PostType, postSlice) {
		return custom.ParamErrorReturn(postType)
	}

	if !utils.StrInSlice(in.Status, constant.TaskStatusSlice) {
		return custom.ParamErrorReturn(status)
	}

	if in.Remark != "" && len(in.Remark) > 1024 {
		return custom.ParamErrorReturn(remark)
	}

	_, err := time.ParseInLocation(constant.TIMELAYOUT, in.NextRuntime, utils.DefaultLocation)
	if err != nil {
		return custom.ParamErrorReturn(nextRuntime)
	}

	if in.Protocol == constant.HTTPJOB {
		if !strings.HasPrefix(in.Command, "http") && !strings.HasPrefix(in.Command, "https") {
			in.Command = fmt.Sprintf("%s%s", "http://", in.Command)
		}
	}

	if in.Headers != "" {
		if len(in.Headers) > 2048 {
			return custom.ParamErrorReturn(headers)
		}

		var m http.Header
		if err := json.Unmarshal([]byte(in.Headers), &m); err != nil {
			return custom.ParamErrorReturn(headers)
		}

	}

	return nil
}
