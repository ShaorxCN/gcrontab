package tasklog

import (
	"strconv"
	"testing"
	"time"
)

type Test struct {
	Name string
}

func TestAppend(t *testing.T) {
	var ok bool
	var te *Test

	testChan := make(chan *Test, 100)
	res := make([]*Test, 0, 10)
	go func() {
		ticket := time.NewTicker(1 * time.Second)
		i := 0
		for {
			tte := new(Test)
			tte.Name = strconv.Itoa(i)
			testChan <- tte
			<-ticket.C
			i++

			if i == 10 {
				close(testChan)
				break
			}

		}
	}()
	for {
		if te, ok = <-testChan; ok {
			res = append(res, te)
		} else {
			break
		}
	}

	for _, v := range res {
		t.Log(*v)
	}
}
