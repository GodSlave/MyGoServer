package user

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/utils"
	"github.com/go-xorm/xorm"
	"time"
	"github.com/GodSlave/MyGoServer/base"
)

type ModuleUser struct {
	basemodule.BaseModule
	redisPool *redis.Pool
	sqlEngine *xorm.Engine
	app       module.App
}

type BaseUser struct {
	Name     string `xorm:"unique index" json:"-"`
	Phone    string `xorm:"unique index"`
	Password string
	Id       string    `xorm:"pk"`
	CreateAt time.Time `xorm:"created"`
}

var Module = func() module.Module {
	newGate := new(ModuleUser)

	return newGate
}

func (m *ModuleUser) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	url := settings.Settings["redis"].(string)
	m.redisPool = utils.GetRedisFactory().GetPool(url)
	m.sqlEngine = app.GetSqlEngine()
	m.GetServer().RegisterGO("Login", m.Login)
	m.app = app

	var user = &BaseUser{}
	m.sqlEngine.Sync(user)
}

func (m *ModuleUser) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "User"
}
func (m *ModuleUser) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

func (m *ModuleUser) Run(closeSig chan bool) {

	<-closeSig
	m.redisPool.Close()
}

func (m *ModuleUser) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}

type LoginForm struct {
	Name     string    `json:"name"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	VerifyCode string `json:"verifyCode"`
}

func (m *ModuleUser) Login(SessionId string, form *LoginForm) (result string, err *base.ErrorCode) {
	user := new(BaseUser)
	has, err1 := m.sqlEngine.Where("name=?", form.Name).Get(user)
	if err1 == nil && has {

	} else {
		return "", base.ErrLoginFail
	}

	return "success", base.ErrNil
}

func Regiester(SessionId string, form RegisterForm) (result string, err string) {

	return "success", ""
}
