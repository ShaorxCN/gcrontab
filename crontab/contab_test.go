package crontab

import (
	"fmt"
	"testing"
	"time"
)

func TestCrontab(t *testing.T) {
	c := new(CrontabConfig)
	c.Init()
	t.Log(ts)
}

type test struct {
	name string
	age  int
}

func pNameAndAgePlus(t *test) {
	fmt.Println(t.name)
	t.age++
}

func TestClosureIn122(t *testing.T) {

	ts := []*test{{"1", 1}, {"2", 2}, {"3", 3}, {"4", 4}, {"5", 5}}

	for _, v := range ts {
		go func() {
			pNameAndAgePlus(v)
		}()
	}

	time.Sleep(3 * time.Second)

	for _, v := range ts {
		t.Log(*v)
	}
}
