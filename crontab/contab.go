package crontab

import (
	"context"
	"encoding/json"
	"fmt"
	"gcrontab/constant"
	"gcrontab/email"
	"gcrontab/entity/task"
	tasklog "gcrontab/entity/task_log"
	"gcrontab/rep/requestmodel"
	taskRep "gcrontab/rep/task"
	taskLogRep "gcrontab/rep/task_log"
	"log"
	"net/http"

	"gcrontab/custom"
	"gcrontab/model"
	"gcrontab/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"os"
	"runtime"
	"sync"
	"time"
)

var (
	ts   *taskScheduler
	todo []*task.Task
)

type taskScheduler struct {
	ctx    context.Context
	cancel context.CancelFunc
	// 参数
	MaxGoroutine chan int
	//  单位 ms
	DBScanInterval  int
	MemScanInterval int
	wg              sync.WaitGroup
	exit            chan struct{}
	// 更新之后需要在该周期内执行的程序channel
	updateTaskChan chan *task.Task
	// 更新继续执行的任务数组
	updateExecTask []*task.Task
	lock           sync.Mutex
	// 记录todo中任务uuid和index的映射
	taskUUIDMap map[string]int
}

// Stop 关闭调度
func (ts *taskScheduler) Stop() {
	close(ts.exit)
	close(ts.updateTaskChan)
	ts.cancel()
	ts.wg.Wait()
}

func (ts *taskScheduler) Start() {
	ts.wg.Add(1)
	go ts.dealUpdateTask()
	go ts.schedulerStart()
}

func getLockInMap(t *task.Task) bool {
	return utils.RegisterEntityInRedis(t, constant.Host, t.Expired_time/1000)

}

func unLockInMap(t *task.Task) error {
	return utils.UnregisterEntityInRedis(t)
}

// ExecImmediately 立即执行
func ExecImmediately(t *task.Task, operator string) error {
	hostName, err := os.Hostname()
	if err != nil {
		logrus.Errorf("get hostname failed:%v", err)
		return err
	}
	// 先创建日志记录
	now := utils.Now()
	tl := &tasklog.TaskLog{
		TimeStamp: now.UnixNano(),
		Status:    constant.STATUSPROCE,
		TaskName:  t.Name,
		TaskID:    t.ID.GetIDValue(),
		Command:   t.Command,
		StartTime: now.String(),
		User:      operator,
		Host:      hostName,
	}

	ts.handler(t, tl)
	return nil
}

func (ts *taskScheduler) dealUpdateTask() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("handler updateTaskChan panic:%s", string(buf[:n]))
		}
	}()

	var ok bool
	var t *task.Task
	for {
		if t, ok = <-ts.updateTaskChan; ok {
			ts.appendupdateExecTask(t)
		} else {
			return
		}
	}
}

func (ts *taskScheduler) appendupdateExecTask(t *task.Task) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.updateExecTask = append(ts.updateExecTask, t)
}

func (ts *taskScheduler) schedulerStart() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("handler crontab panic:%s", string(buf[:n]))
		}
	}()
	var deadline time.Time
	var err error
	tickerInDB := time.NewTicker(time.Duration(ts.DBScanInterval) * time.Millisecond)
	tickerInMem := time.NewTicker(1 * time.Second)
	// idleTimeDuration := time.After(1 * time.Second)
	for {
		select {
		case <-ts.exit:
			for {
				done := len(ts.MaxGoroutine)
				if done != 0 {
					logger.WithTime(utils.Now()).Infof("等待剩余goroutine运行结束,运行中:%d", done)
					time.Sleep(time.Second * 1)
					continue
				}
				ts.wg.Done()
				return
			}
		default:
		}
		now := utils.Now()
		deadline = utils.Now().Add(time.Duration(ts.DBScanInterval) * time.Millisecond)
		taskrep := taskRep.NewTaskRep(model.DB())
		todo, err = taskrep.FindActiveTasks(deadline)
		if err != nil {
			// TODO:是否需要通知
			logger.WithTime(utils.Now()).Errorf("find active tasks failed:%v", err)
		}

	memScan:
		for {
			if deadline.Before(utils.Now()) {
				break
			}

			select {
			case <-ts.ctx.Done():
				break memScan
			default:
			}
			// 1.22 新的loopvar 每次都会重新声明定义te 不需要重新传
			for index, te := range todo {
				if te == nil {
					continue
				}

				ts.MaxGoroutine <- 1
				go func(index int, t *task.Task) {
					defer func() {
						<-ts.MaxGoroutine
						if err := recover(); err != nil {
							buf := make([]byte, 2048)
							n := runtime.Stack(buf, false)
							logger.WithTime(utils.Now()).Errorf("handler task[%s] panic:%s", t.ID, string(buf[:n]))
						}
					}()
					if !getLockInMap(t) {
						logger.WithTime(utils.Now()).Warnf("[%s]TaskName[%s] get Lock failed", t.ID.GetIDValue(), t.Name)
						return
					}

					if nt, needRun := doubleCheck(now, t.ID.GetIDValue()); !needRun {
						logger.WithTime(utils.Now()).Warnf("[%s]TaskName[%s] double check failed", t.ID.GetIDValue(), t.Name)
						err := unLockInMap(t)
						if err != nil {
							logger.WithTime(utils.Now()).Errorf("unlock entity failed:%v entity:%v", err, t.ID)
						}
						if utils.IsBeforeOrEq(nt.NextRuntimeUse, deadline) {
							ts.updateTaskChan <- t
						} else {
							todo[index] = nil
						}
						return
					} else {
						// 选择使用最新查询到的或者还是使用第一次查询到的
						t = nt
					}

					defer func() {
						err := unLockInMap(t)
						if err != nil {
							logger.WithTime(utils.Now()).Errorf("unlock entity failed:%v entity:%v", err, t.ID)
						}
					}()
					execTime := utils.Now()
					// 先记录日志
					tl, err := saveTaskLog(t, execTime)
					if err != nil {
						logger.WithTime(utils.Now()).Errorf("sava task[%v] start exec log error:%v", t.ID, err)
						return
					}
					logger.WithTime(utils.Now()).Infof("[%s]exec...", t.ID)
					remove := ts.handler(t, tl)

					if remove {
						todo[index] = nil
					} else {
						nextTime, err := updateTask4Next(t, execTime)
						if err != nil {
							logger.WithTime(utils.Now()).Errorf("update task[%s] to next failed :%v", t.ID, err)
							return
						}

						// 这里可以直接修改 不用过channel
						if utils.IsBeforeOrEq(nextTime, deadline) {
							ts.updateTaskChan <- t
						} else {
							todo[index] = nil
						}
						logger.WithTime(utils.Now()).Infof("[%s]exec end...", t.ID)
					}
				}(index, te)
			}

			<-tickerInMem.C
			ts.dealUpdateTaskSlice()
		}

		select {
		case <-tickerInDB.C:
		case <-ts.ctx.Done():
			logrus.Infoln("cancel crontab.....")
		}
	}
}

func (ts *taskScheduler) dealUpdateTaskSlice() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	for _, v := range ts.updateExecTask {
		if v != nil {
			if index, ok := ts.taskUUIDMap[v.ID.GetIDValue()]; ok {
				todo[index] = v
			} else {
				todo = append(todo, v)
				ts.taskUUIDMap[v.ID.GetIDValue()] = len(todo) - 1
			}
		}
	}
}

func doubleCheck(now time.Time, id string) (*task.Task, bool) {
	taskRep := taskRep.NewTaskRep(nil)
	tnew, err := taskRep.FindTaskByID(uuid.MustParse(id))
	if err != nil {
		logger.WithTime(utils.Now()).Errorf("find task by id[%s] failed:%v", id, err)
		return tnew, false
	}

	return tnew, utils.IsBeforeOrEq(tnew.NextRuntimeUse, now)
}

func (ts *taskScheduler) handler(t *task.Task, tl *tasklog.TaskLog) bool {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("handler task panic:%s", string(buf[:n]))
		}
	}()
	switch t.Protocol {
	case constant.HTTPJOB:
		httpHandler(t, tl)
	default:
		logger.WithTime(utils.Now()).Errorf("not support this type of job[%s]", t.Command)
		protocolFailHandler(t, tl)
		return true

	}

	return false
}

func saveTaskLog(t *task.Task, tm time.Time) (tl *tasklog.TaskLog, err error) {
	defer func() {
		if errPanic := recover(); errPanic != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("save TaskLog to DB panic:%v,%s", errPanic, string(buf[:n]))
			err = custom.ErrorInternalServerError
		}
	}()

	tl = new(tasklog.TaskLog)
	tl.TimeStamp = tm.UnixNano()
	tl.TaskName = t.Name
	tl.TaskID = t.ID.GetIDValue()
	tl.Status = constant.STATUSPROCE
	tl.Command = t.Command
	tl.StartTime = tm.String()
	tl.StartTimeT = tm
	hostName, _ := os.Hostname()
	tl.Host = hostName
	logRep := taskLogRep.NewTaskLogRep(model.DB())
	err = logRep.SaveTaskLog(tl)
	if err != nil {
		logger.WithTime(utils.Now()).Errorf("save task_log failed:%v", err)
	}

	return
}

// TODO:协议不支持的不落库 只发邮件以及记录日志 ?是否需要重复记录不支持的任务类型
func protocolFailHandler(t *task.Task, tl *tasklog.TaskLog) {
	copy := *t

	go func(t *task.Task) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				log.Printf("%v handler task panic:%s", time.Now(), string(buf[:n]))
			}
		}()

		// 统一设置的通知email
		emails, err := model.FindEmails(t.Creater)
		if err != nil {
			logrus.Errorf("taskID[%s] find  email addresses failed:%v", t.ID, err)
			return
		}
		err = email.SendCrontabAlert(tl.ResultCode, tl.Result, t, tl.TimeStamp, emails)
		if err != nil {
			logrus.Errorf("taskID[%s] send crontab alert email[To:%s] failed:%v", t.ID, emails, err)
		}

	}(&copy)

	res := &ResponseWrapper{StatusCode: -1, Body: custom.ErrorUnSupportTaskProtocol.Error(), Header: make(http.Header)}
	res.End = utils.Now()
	updateTaskLog(res, tl)
}

func httpHandler(t *task.Task, tl *tasklog.TaskLog) {
	var res *ResponseWrapper
	h := make(http.Header)

	if t.Headers != "" {
		err := json.Unmarshal([]byte(t.Headers), &h)
		if err != nil {
			logger.WithTime(utils.Now()).Errorf("task header unmarshal failed:%v,data:[%s]", err, t.Headers)
		}
	}

	switch t.HTTPMethod {
	case constant.HTTPMETHODGET:
		res = Get(t.Command, t.Expired_time, &h)
	case constant.HTTPMETHODPOST:
		switch t.PostType {
		case constant.POSTJSON:
			res = PostJSON(t.Command, t.Param, t.Expired_time, &h)
		case constant.POSTFORM:
			res = PostForm(t.Command, t.Param, t.Expired_time, &h)
		default:
			res = &ResponseWrapper{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       fmt.Sprintf("unsupport post type:%s", t.PostType),
				Start:      utils.Now(),
				End:        utils.Now(),
			}
		}
	}

	if res.StatusCode != http.StatusOK {
		emails, err := model.FindEmails(t.Creater)
		if err != nil {
			logger.WithTime(utils.Now()).Errorf("taskID[%s] find  email addresses failed:%v", t.ID, err)
			return
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					log.Printf("%v handler task panic:%s", time.Now(), string(buf[:n]))
				}
			}()

			err = email.SendCrontabAlert(res.StatusCode, res.Body, t, tl.TimeStamp, emails)
			if err != nil {
				logger.WithTime(utils.Now()).Errorf("taskID[%s] send crontab alert email[To:%s] failed:%v", t.ID, emails, err)
			}
		}()
	}
	updateTaskLog(res, tl)
}

func updateTaskLog(res *ResponseWrapper, tl *tasklog.TaskLog) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("save TaskLog to DB panic:%v,%s", err, string(buf[:n]))
		}
	}()

	tl.ResultCode = res.StatusCode
	tl.Status = constant.STATUSSUCC
	if res.StatusCode != 200 {
		tl.Status = constant.STATUSFAIL
	}
	tl.Result = res.Body
	tl.TotalTime = res.End.Sub(tl.StartTimeT).Nanoseconds() / 1e6
	tl.EndTime = res.End.String()
	logrep := taskLogRep.NewTaskLogRep(model.DB())
	err := logrep.UpdateTaskLog(tl)
	if err != nil {
		logger.WithTime(utils.Now()).Errorf("update task_log failed:%v", err)
	}
}

func updateTask4Next(t *task.Task, exec time.Time) (time.Time, error) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("modify task to nextExec panic:%s", string(buf[:n]))
		}
	}()

	param := &requestmodel.ModifyTask{}

	param.LastRuntimeUse = exec
	param.NextRuntimeUse = utils.GetNextTimeAfterNow(t.NextRuntimeUse, t.IntervalDuration, t.UnitOfInterval).In(utils.DefaultLocation)
	param.NextRuntime = param.NextRuntimeUse.Format(constant.TIMELAYOUT)
	param.UpdateFlag = 0

	taskID, err := uuid.Parse(t.ID.GetIDValue())
	if err != nil {
		return param.NextRuntimeUse, err
	}
	t.LastRuntimeUse = exec
	t.NextRuntimeUse = param.NextRuntimeUse
	taskRep := taskRep.NewTaskRep(nil)

	return param.NextRuntimeUse, taskRep.ModifyTaskTimeByID(taskID, param)
}
