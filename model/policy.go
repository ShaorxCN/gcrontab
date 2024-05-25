package model

type DBPolicy struct {
	ID     int
	Ptype  string `gorm:"not null;type:varchar(18);default:'p'"`
	Role   string `gorm:"not null;type:varchar(1024)"`
	Path   string `gorm:"not null;type:varchar(1024)"`
	Method string `gorm:"not null;type:varchar(1024)"`
}

var policyTableName = new(DBPolicy).TableName()

func (DBPolicy) TableName() string {
	return "tbl_policy"
}

func FindAllPolicy() ([]*DBPolicy, error) {
	var ps []*DBPolicy
	db := DB().Table(policyTableName)
	err := db.Find(&ps).Error
	return ps, err
}

func SavePolicy(p *DBPolicy) error {
	return db.Create(p).Error
}

func FindPolicyByID(id int) (*DBPolicy, error) {
	p := new(DBPolicy)
	p.ID = id
	db := DB()
	err := db.Model(p).First(p).Error

	return p, err
}
