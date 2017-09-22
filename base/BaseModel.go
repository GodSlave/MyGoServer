package base

var SESSION_PERFIX = "session"
var TOKEN_PERFIX = "token"
var REFRESH_TOKEN_PERFIX = "rToken"

// DB model
type BaseUser struct {
	// @inject_tag: xorm:"unique index notnull"
	Name      string `xorm:"unique index notnull"`
	Phone     string `xorm:"index"`
	Password  string `xorm:"notnull"`
	UserID    string `xorm:"unique index notnull"`
	Id        int64
	CreatTime int64
}