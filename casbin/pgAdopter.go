package casbin

import (
	"bufio"
	"bytes"
	"errors"
	dbmodel "gcrontab/model"
	"strings"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
	"github.com/jinzhu/gorm"
)

// call LoadPolicy after save a policy record to db ,if you want the policy  saved effects immediately

const (
	m = `[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = r.sub == p.sub && regexMatch(r.obj, p.obj) && regexMatch(r.act, p.act)||r.sub=="admin"`

	ps = `p, admin, .*, (GET)|(PUT)|(POST)|(DELETE)
	p, user, /tasks.*, GET
	p, user, /users/.+, (GET)|(PUT)
	p, user, /taskLogs.*, GET
	p, taskAdmin, /tasks.*, (POST)|(GET)|(PUT)|(DELETE)
	p, taskAdmin, /taskLogs.*, GET
	p, taskAdmin, /users/.+, (GET)|(PUT)
	p, user, /view/taskLogs, GET
	p, taskAdmin, /view/taskLogs, GET`
)

var (
	// CasEnforecer crontab 权限管理示例
	CasEnforecer *casbin.Enforcer
)

// Init casbin 初始化
func Init() error {
	scanner := bufio.NewScanner(bytes.NewBuffer([]byte(ps)))
	i := -1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		cols := strings.Split(line, ",")
		p := new(dbmodel.DBPolicy)
		p.ID = i
		p.Ptype = strings.TrimSpace(cols[0])
		p.Role = strings.TrimSpace(cols[1])
		p.Path = strings.TrimSpace(cols[2])
		p.Method = strings.TrimSpace(cols[3])

		_, err := dbmodel.FindPolicyByID(i)
		i--
		if err == nil {
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		err = dbmodel.SavePolicy(p)
		if err != nil {
			return err
		}

	}
	pg := new(PgAdopter4Crontab)
	mm := casbin.NewModel(m)
	CasEnforecer = casbin.NewEnforcer(mm, pg)

	return nil
}

// PgAdopter4Crontab 给crontab使用的policy adopter
type PgAdopter4Crontab struct {
}

// LoadPolicy loads all policy rules from the storage.
func (a *PgAdopter4Crontab) LoadPolicy(m model.Model) error {
	ps, err := dbmodel.FindAllPolicy()
	if err != nil {
		return err
	}

	var s []string
	for _, v := range ps {
		s = make([]string, 3)
		s[0] = v.Role
		s[1] = v.Path
		s[2] = v.Method
		m[v.Ptype][v.Ptype].Policy = append(m[v.Ptype][v.Ptype].Policy, s)
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *PgAdopter4Crontab) SavePolicy(m model.Model) error {
	var p *dbmodel.DBPolicy
	var err error
	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			p = new(dbmodel.DBPolicy)
			p.Ptype = ptype
			if len(rule) > 0 {
				p.Role = rule[0]
			}
			if len(rule) > 1 {
				p.Path = rule[1]
			}
			if len(rule) > 2 {
				p.Method = rule[2]
			}

			err = dbmodel.SavePolicy(p)
			if err != nil {
				return err
			}
		}

	}

	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *PgAdopter4Crontab) AddPolicy(sec, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *PgAdopter4Crontab) RemovePolicy(sec, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *PgAdopter4Crontab) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
