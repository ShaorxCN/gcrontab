package crontab

import (
	"testing"
)

func TestCrontab(t *testing.T) {
	c := new(CrontabConfig)
	c.Init()
	t.Log(ts)
}
