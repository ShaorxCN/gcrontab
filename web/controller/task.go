package controller

import (
	"github.com/gin-gonic/gin"
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

}

// RunTask 立即执行一条任务。 只更新日志
func (s Task) RunTask(ctx *gin.Context) {

}

// ModifyTask 修改一条task 添加标识位 调度任务完成的时候就不会去修改下次执行时间
func (s Task) ModifyTask(ctx *gin.Context) {
}

func (s Task) DeleteTask(ctx *gin.Context) {

}

// FindTaskByID 根据id 查找任务
func (s Task) FindTaskByID(ctx *gin.Context) {

}

func (s Task) FindTasks(ctx *gin.Context) {

}

func (s Task) CreateTaskByAPIV1(ctx *gin.Context) {

}
