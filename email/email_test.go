package email

import (
	"gcrontab/entity/task"
	"gcrontab/interface/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestSendCrontabAlert(t *testing.T) {
	c := Config{
		Addr:     "smtp.exmail.qq.com:465",
		Host:     "smtp.exmail.qq.com",
		PassWord: "xxxxxxx",
		User:     "xxxx@xxxx.com",
	}

	err := c.Init()
	if err != nil {
		logrus.Error("init error: ", err)
		return
	}

	taskTest := &task.Task{Name: "test"}
	taskTest.ID = entity.NewEntityKey(uuid.New().String(), task.TaskEntityType)

	err = SendCrontabAlert(404, "this is error", taskTest, 123, []string{"xxxx@xxxxx.com"})
	if err != nil {
		logrus.Error(err)
	}
}
