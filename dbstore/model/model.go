package model

type Admin struct {
	ID       uint64 `xorm:"pk autoincr 'id'"`
	Login    string `xorm:"login"`
	Password []byte `xorm:"password"`
}

func (Admin) TableName() string {
	return "admins"
}

type Project struct {
	ID   string `xorm:"pk 'id'"`
	Name string `xorm:"name"`
}

func (Project) TableName() string {
	return "projects"
}

type User struct {
	ID       uint64 `xorm:"pk autoincr 'id'"`
	Email    string `xorm:"email"`
	Password []byte `xorm:"password"`
	Name     string `xorm:"name"`
}

func (User) TableName() string {
	return "users"
}

type UserAddress struct {
	UserID  uint64 `xorm:"userid"`
	Coin    string `xorm:"coin"`
	Address string `xorm:"address"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
