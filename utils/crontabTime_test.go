package utils

import (
	"testing"
	"time"
)

func TestIsBeforeOrEq(t *testing.T) {
	st := time.Now()
	t_eq := st
	t_over := st.Add(time.Second)

	t_before := st.Add(-1 * time.Second)

	tests := []struct {
		name string
		arg  time.Time
		want bool
	}{
		{"before", t_before, true}, {"eq", t_eq, true}, {"after", t_over, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if res := IsBeforeOrEq(test.arg, st); res != test.want {
				t.Errorf("IsBeforeOrEq()=%v,want %v", res, test.want)
			}
		})
	}
}
