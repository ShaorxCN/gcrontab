package custom

import (
	"errors"
	"fmt"
)

const (
	StatusOK                   = "200"
	StatusAccepted             = "202"
	StatusNonAuthoritativeInfo = "203"
	StatusBadRequest           = "400"
	StatusNotFound             = "404"
	StatusForbidden            = "403"
	StatusLocked               = "423"
	StatusFailedDependency     = "424"
	InternalServerError        = "500"
	ParamError                 = "461"
	RecordExist                = "462"
	APICheckSignError          = "100"
)

var (
	ErroEntityKeyParam         = errors.New("neither idValue or tpValue can be empty")
	ErrorInvalideAccessToken   = errors.New("登陆状态失效")
	ErrorSaveToDBFailed        = errors.New("操作数据库/redis失败")
	ErrorInvalideRequest       = errors.New("请求格式错误")
	ErrorCalculateTimeFailed   = errors.New("计算下次执行时间出错")
	ErrorEntityLocked          = errors.New("系统忙 请稍后重试")
	ErrorRecordNotFound        = errors.New("未找到记录")
	ErrorLoginFailed           = errors.New("用户名或密码错误")
	ErrorInternalServerError   = errors.New("系统内部错误")
	ErrorRecordExist           = errors.New("资源已经存在")
	ErrorOriginPasswordError   = errors.New("原密码错误")
	ErrorUnSupportTaskProtocol = errors.New("不支持的任务类型")
	ErrorForbidden             = errors.New("权限不足")
	ErrorSignError             = errors.New("签名错误")
	ErrorRunTaskError          = errors.New("执行任务失败")
)

// ParamErrorReturn 字段校验错误
func ParamErrorReturn(field string) error {
	return fmt.Errorf("[%s]字段错误", field)
}
