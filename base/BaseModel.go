package base

var SESSION_PERFIX = "session"
var ID_SESSION_PREFIX = "USession"
var TOKEN_PERFIX = "token"
var ID_TOKEN_PERFIX = "UToken"
var REFRESH_TOKEN_PERFIX = "rToken"
var ID_REFRESH_TOKEN_PREFIX = "URTOKEN"

// DB model
type BaseUser struct {
	// @inject_tag: xorm:"unique index notnull"
	Name      string `xorm:"unique index notnull"`
	Phone     string `xorm:"index"`
	Password  string `xorm:"notnull"`
	UserID    string `xorm:"unique index notnull"`
	Age       int32  `xorm:"notnull"`
	Id        int64
	CreatTime int64
}

type PushItem struct {
	Module   byte
	PushType byte
	Content  []byte
}
