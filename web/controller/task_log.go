package controller

import (
	"gcrontab/custom"
	"gcrontab/rep/requestmodel"
	"gcrontab/web/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User 用来实现用户的 rest 接口。
type TaskLog struct{}

// AddTaskRouter 注册用户 router。
func AddTaskLogRouter(e *gin.Engine) {
	e.GET("/taskLogs/:taskID", TaskLog{}.FindTaskLogByTaskID)
	// e.GET("/taskLogs", TaskLog{}.FindTaskLogs)
	// e.GET("/view/taskLogs", TaskLog{}.FindTaskLogView)
}

// FindTaskLogByTaskID 根据任务id查找任务日志接口。
func (s TaskLog) FindTaskLogByTaskID(ctx *gin.Context) {
	params := ctx.Value("params").(requestmodel.Params)
	taskID, err := uuid.Parse(ctx.Param("taskID"))
	if err != nil {
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ParamErrorReturn("taskID").Error()))
		return
	}
	params.ID = taskID

}
