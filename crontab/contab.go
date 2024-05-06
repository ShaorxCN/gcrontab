package gocron

/**
 * 任务完成或者超时才记录日志
 * 20190805 执行前记录日志  完成后更新
 */

import (
	"gcrontab/constant"

	"gcrontab/custom"
	"gcrontab/model"
	"gcrontab/utils"

	"github.com/google/uuid"

	"os"
	"runtime"
	"sync"
	"time"
)

var (
	ts *taskScheduler
	// imStop      chan int
	// sStop       chan int

	// TaskChannel 任务channel
	TaskChannel chan *model.DBTask
	// PeoPlusKey 调用P+的key
	PeoPlusKey string
)

type taskScheduler struct {
	// 参数
	MaxGoroutine chan int
	//  单位 ms
	ScanInterval int
	wg           sync.WaitGroup
	exit         chan struct{}
}

// Stop 关闭调度
func (ts *taskScheduler) Stop() {
	close(ts.exit)
	ts.wg.Wait()
}

func (ts *taskScheduler) Start() {
	ts.wg.Add(1)
	go ts.schedulerStart()
}

// ExecImmediately 立即执行
func ExecImmediately(t *model.DBTask, tl *model.DBTaskLog) {
	ts.handler(t, tl)
}

// TODO: 改成例如2h内需要执行的甚至当天需要执行  然后根据时间排序？ 或者all go 然后time.after? add or update ？
// TODO: 2h扫描有一次 然后内存保持维护排序 每秒扫描检查是否需要执行？如果nexttime 计算下如果还是这个周期内则修改完后继续加入队列 否则只落到db
// 因为放弃使用所有go出去time.after 阻塞 而是每秒扫描时间排序  数据结构 map[id]entity 还有个[]time.Time 实现time sort
func (ts *taskScheduler) schedulerStart() {

	ticker := time.NewTicker(time.Duration(ts.ScanInterval) * time.Millisecond)
	// idleTimeDuration := time.After(1 * time.Second)
	for {
		now := utils.Now()
		tasks, err := model.FindActiveTasks(now)
		if err != nil {
			// TODO:notify
			logger.WithTime(utils.Now()).Errorf("find active tasks failed:%v", err)
		}
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

		for _, task := range tasks {
			ts.MaxGoroutine <- 1
			go func(t *model.DBTask) {
				defer func() {
					<-ts.MaxGoroutine
					if err := recover(); err != nil {
						buf := make([]byte, 2048)
						n := runtime.Stack(buf, false)
						logger.WithTime(utils.Now()).Errorf("handler task[%s] panic:%s", t.ID, string(buf[:n]))
					}
				}()
				if !getLockInMap(t) {
					logger.WithTime(utils.Now()).Warnf("[%s]TaskName[%s] get Lock failed", t.ID.String(), t.Name)
					return
				}

				if !doubleCheck(now, t.ID) {
					logger.WithTime(utils.Now()).Warnf("[%s]TaskName[%s] double check failed", t.ID.String(), t.Name)
					err := unLockInMap(t)
					if err != nil {
						logger.WithTime(utils.Now()).Errorf("unlock entity failed:%v entity:%v", err, t.ID)
					}
					return
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
				ts.handler(t, tl)

			}(task)
		}
		<-ticker.C
	}
}

func doubleCheck(now time.Time, id uuid.UUID) bool {
	tnew, err := model.FindTaskByID(id)
	if err != nil {
		logger.WithTime(utils.Now()).Errorf("find task by id[%s] failed:%v", id.String(), err)
		return false
	}
	return tnew.NextRuntime.Before(now)
}

func (ts *taskScheduler) handler(t *model.DBTask, tl *model.DBTaskLog) {
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
	// case constant.EXECJON:
	// 	fallthrough
	default:
		logger.WithTime(utils.Now()).Errorf("not support this type of job[%s]", t.Command)
		failHandler(tl, t)

	}

}

func saveTaskLog(t *model.DBTask, tm time.Time) (tl *model.DBTaskLog, err error) {
	defer func() {
		if errPanic := recover(); errPanic != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.WithTime(utils.Now()).Errorf("save TaskLog to DB panic:%v,%s", errPanic, string(buf[:n]))
			err = custom.ErrorInternalServerError
		}
	}()

	tl = new(model.DBTaskLog)
	tl.TimeStamp = tm.UnixNano()
	tl.TaskName = t.Name
	tl.TaskID = t.ID
	tl.Status = constant.STATUSPROCE
	tl.Command = t.Command
	tl.StartTime = tm
	hostName, _ := os.Hostname()
	tl.Host = hostName
	err = model.InsertTaskLog(tl)
	if err != nil {
		logger.WithTime(utils.Now()).Errorf("save task_log failed:%v", err)
	}

	return
}
