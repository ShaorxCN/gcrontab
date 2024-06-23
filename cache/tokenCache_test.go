package cache

import "testing"

func TestTokenCahce(t *testing.T) {
	SaltCacheInit(30)

	type setting struct {
		uid   string
		token string
		salt  string
	}

	type wants struct {
		salt string
		ok   bool
	}

	tests := []struct {
		name string
		set  *setting
		arg  struct {
			uid   string
			token string
		}
		want wants
	}{{name: "existswithoutupdate", set: &setting{"exists", "extoken", "exsalt"}, arg: struct {
		uid   string
		token string
	}{"exists", "extoken"}, want: wants{"exsalt", true}}, {name: "exists2", set: &setting{"exists", "extoken2", "exsalt2"}, arg: struct {
		uid   string
		token string
	}{"exists", "extoken2"}, want: wants{"exsalt2", true}}, {name: "existswithupdate", set: &setting{"exists", "extoken", "exsalt3"}, arg: struct {
		uid   string
		token string
	}{"exists", "extoken"}, want: wants{"exsalt3", true}}, {"none", nil, struct {
		uid   string
		token string
	}{"none", "notoken"}, wants{"", false}}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.set != nil {
				SetSalt(test.set.uid, test.set.token, test.set.salt)
			}

			get, ok := GetSaltByUIDAndToken(test.arg.uid, test.arg.token)
			if ok != test.want.ok || get != test.want.salt {
				t.Errorf("tokenCache_test failed,want:%v,get:{%v,%v}", test.want, get, ok)
			}
		})
	}

	RemoveByUIDAndToken(tests[1].arg.uid, tests[1].arg.token)
	_, ok := GetSaltByUIDAndToken(tests[1].arg.uid, tests[1].arg.token)

	if ok {
		t.Error("get true after del,error")
	}

	RemoveUIDSalt(tests[0].arg.uid)

	_, ok = GetSaltByUIDAndToken(tests[0].arg.uid, tests[0].arg.token)

	if ok {
		t.Error("get true after del,error")
	}

}
