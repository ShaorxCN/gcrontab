package main

import (
	"errors"
	"fmt"
	"gcrontab/config"
	"gcrontab/crontab"
	"gcrontab/model"
	"gcrontab/redis"
	"gcrontab/utils"
	"gcrontab/web"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/sirupsen/logrus"
)

func startPprofServe(port int) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("SYSTEM ACTION PANIC: %v", r)
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				logrus.Error(string(buf[:n]))
			}
		}()

		logrus.Info(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()
}

func main() {
	errChan := make(chan error)

	cm := config.DefaultManagement
	cm.Register("db_config", &model.DbConfig{})
	cm.Register("gin_config", &web.GinConfig{})
	cm.Register("crontab_config", &crontab.CrontabConfig{})
	cm.Register("redis_config", &redis.Config{})
	// cm.Register("email_config", &email.Config{})
	err := cm.Init()
	if err != nil {
		logrus.Fatalf("init config failed, err: %v", err)
	}

	statePort := cm.Configs["gin_config"].(*web.GinConfig).StatePort
	startPprofServe(utils.If(statePort == 0, 3998, statePort).(int))

	err = cm.Configs["redis_config"].Init()
	if err != nil {
		logrus.Fatalf("could not connect to the redis:%s", err.Error())
	}

	err = cm.Configs["db_config"].Init()
	if err != nil {
		logrus.Fatalf("could not connect to the db:%s", err.Error())
	}

	// err = cm.Configs["email_config"].Init()
	// if err != nil {
	// 	logrus.Fatalf("could not connect to the email:%s", err.Error())
	// }
	err = model.AutoMigrate()
	if err != nil {
		logrus.Fatalf("create tables failed:%s", err.Error())
	}
	go (func() {

		errChan <- cm.Configs["gin_config"].Init()

	})()
	go func() { errChan <- cm.Configs["crontab_config"].Init() }()

	stopFalg := make(chan os.Signal, 1)
	signal.Notify(stopFalg, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopFalg
		logrus.StandardLogger().WithTime(utils.Now()).Warn("ready to stop process...")
		cm.Stop()
		errChan <- errors.New("stop process")
	}()

	for err := range errChan {
		if err != nil {
			logrus.StandardLogger().WithTime(utils.Now()).Fatalf(err.Error())
		}
	}
}
