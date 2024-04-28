package security

import "testing"

func TestHashSha256(t *testing.T) {
	origin := "this is for test"
	res := "038a1860a4192f482250af54737a0f2e5f5b90763965a095144e91ed15077977"

	if HashSha256(origin) != res {
		t.Fail()
	} else {
		t.Log("pass")
	}
}
