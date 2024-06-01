package response

import (
	tasklog "gcrontab/entity/task_log"
)

// BaseResponse web 基础应答
type BaseResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"subCode,omitempty"`
	SubMsg  string `json:"subMsg,omitempty"`
}

// APIResponse api response
type APIResponse struct {
	Code string `json:"subCode"`
	Msg  string `json:"subMsg"`
}

// NewSuccessBaseResponse 成功web res
func NewSuccessBaseResponse() *BaseResponse {
	return &BaseResponse{
		Code:    "000",
		Msg:     "success",
		SubCode: "200",
		SubMsg:  "success",
	}
}

// NewBusinessFailedBaseResponse 失败应答构造 服务器ok 但是业务逻辑上讲错误
func NewBusinessFailedBaseResponse(subCode, subMsg string) *BaseResponse {
	return &BaseResponse{
		Code:    "000",
		Msg:     "success",
		SubCode: subCode,
		SubMsg:  subMsg,
	}
}

// NewSystemFailedBaseResponse 系统错误返回
func NewSystemFailedBaseResponse() *BaseResponse {
	return &BaseResponse{
		Code: "500",
		Msg:  "SYSTEM_ERROR",
	}
}

// NewAPIResponse 返回应答结构体
func NewAPIResponse(code, msg string) *APIResponse {
	return &APIResponse{
		Code: code,
		Msg:  msg,
	}
}

// NewAPISuccessResponse 返回成功应答结构体
func NewAPISuccessResponse() *APIResponse {
	return &APIResponse{
		Code: "000",
		Msg:  "success",
	}
}

// FindTaskLogsResponse log 返回结构
type FindTaskLogsResponse struct {
	*BaseResponse
	Count    int64              `json:"count"`
	TaskLogs []*tasklog.TaskLog `json:"taskLogs"`
}
