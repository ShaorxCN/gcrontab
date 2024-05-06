package entity

// BaseEntity 是基本结构体。
type BaseEntity struct {
	ID       string `json:"-" url:"-"`
	CreateAt string `json:"createAt,omitempty" url:"-"`
	UpdateAt string `json:"updateAt,omitempty" url:"-"`
	Creater  string `json:"creater,omitempty" url:"creater,omitempty"`
}
