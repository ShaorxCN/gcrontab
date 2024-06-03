package controller

import (
	"gcrontab/custom"
	"gcrontab/entity/task"
	"gcrontab/service"
	"gcrontab/utils"
	"gcrontab/web/response"
	"gcrontab/web/validate"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Task 用来实现用户的 rest 接口。
type Task struct{}

// AddTaskRouter 注册用户 router。
func AddTaskRouter(e *gin.Engine) {
	e.POST("/tasks", Task{}.CreateTask)
	e.POST("/tasks/:taskID/run", Task{}.RunTask)
	e.PUT("/tasks/:taskID", Task{}.ModifyTask)
	e.DELETE("/tasks/:taskID", Task{}.DeleteTask)
	e.GET("/tasks/:taskID", Task{}.FindTaskByID)
	e.POST("/api/v1/tasks", Task{}.CreateTaskByAPIV1)
}

// CreateTask 创建一个任务。
func (s Task) CreateTask(ctx *gin.Context) {
	in := new(task.Task)
	err := ctx.BindJSON(in)
	if err != nil {
		logrus.WithError(err).Error()
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ErrorInvalideRequest.Error()))
		return
	}

	err = validate.CheckCreateTaskRequest(in)
	if err != nil {
		logrus.Errorf("checkRquestFailed:%v", err)
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.ParamError, err.Error()))
		return
	}
	tasks := []*task.Task{in}
	taskService := service.NewTaskService(utils.NewServiceContext(ctx, nil), nil, tasks)
	err = taskService.CreateTask()

	// 是否需要返回创建的实体 前端以此获取主键方便查询数据做展示
	if err != nil {
		if err == custom.ErrorRecordExist {
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.RecordExist, err.Error()))
			return
		}
		logrus.Errorf("add task  failed:%v", err)
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusFailedDependency, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, response.NewSuccessBaseResponse())
	}

}

// RunTask 立即执行一条任务(补偿用)。 只更新日志 不会修改下次执行时间 不会扰乱原来的周期规律。本身dbscan周期会找到落后的任务
// 并且按照周期规律更新下次执行时间为当前时间后最接近的正常周期时间
func (s Task) RunTask(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("taskID"))
	if err != nil {
		logrus.Errorf("parse uuid failed:%v", err)
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ParamErrorReturn("taskID").Error()))
		return
	}

	taskService := service.NewTaskService(utils.NewServiceContext(ctx, nil), nil, nil)
	err = taskService.RunTask(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusNotFound, custom.ErrorRecordNotFound.Error()))
		} else {
			logrus.Errorf("run task failed:%v", err)
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.InternalServerError, custom.ErrorRunTaskError.Error()))
		}
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.NewSuccessBaseResponse)

}

// ModifyTask 修改一条task 添加标识位 调度任务完成的时候就不会去修改下次执行时间
func (s Task) ModifyTask(ctx *gin.Context) {
}

func (s Task) DeleteTask(ctx *gin.Context) {

}

// FindTaskByID 根据id 查找任务
func (s Task) FindTaskByID(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("taskID"))
	if err != nil {
		logrus.Errorf("parse uuid failed:%v", err)
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ParamErrorReturn("taskID").Error()))
		return
	}

	taskService := service.NewTaskService(utils.NewServiceContext(ctx, nil), nil, nil)
	t, err := taskService.FindTaskByID(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusNotFound, custom.ErrorRecordNotFound.Error()))
		} else {
			logrus.Errorf("find task failed:%v", err)
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusFailedDependency, custom.ErrorEntityLocked.Error()))
		}
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, t)
}

func (s Task) FindTasks(ctx *gin.Context) {

}

func (s Task) CreateTaskByAPIV1(ctx *gin.Context) {

}
