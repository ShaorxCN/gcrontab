package validate

import "testing"

func TestIsValidMD5(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{{"md5", "d41d8cd98f00b204e9800998ecf8427e", true}, {"too short", "notmd5", false}, {"not hex", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if get := IsValidMD5(test.args); get != test.want {
				t.Errorf("IsValidMD5 failed,get()=%v,want %v", get, test.want)
			}
		})
	}
}
