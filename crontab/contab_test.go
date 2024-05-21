package crontab

import (
	"gcrontab/utils"
	"testing"
	"time"
)

func TestCrontab(t *testing.T) {
	c := new(CrontabConfig)
	c.Init()
	t.Log(ts)
}
