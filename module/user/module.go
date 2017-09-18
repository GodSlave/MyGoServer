package user

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/utils"
	"github.com/go-xorm/xorm"
	"github.com/GodSlave/MyGoServer/base"
	"crypto/md5"
	"encoding/hex"
	"github.com/GodSlave/MyGoServer/log"
	"time"
	"github.com/GodSlave/MyGoServer/utils/uuid"
)

type ModuleUser struct {
	basemodule.BaseModule
	redisPool *redis.Pool
	sqlEngine *xorm.Engine
	app       module.App
}

var SESSION_PERFIX = "session"

var Module = func() module.Module {
	newGate := new(ModuleUser)

	return newGate
}

func (m *ModuleUser) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	url := settings.Settings["redis"].(string)
	m.redisPool = utils.GetRedisFactory().GetPool(url)
	m.sqlEngine = app.GetSqlEngine()
	m.GetServer().RegisterGO("Login", 1, m.Login)
	m.GetServer().RegisterGO("Register", 2, m.Regiester)
	m.GetServer().RegisterGO("GetVerifyCode", 3, m.GetVerifyCode)
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

func (m *ModuleUser) Login(SessionId string, form *LoginForm) (result string, err *base.ErrorCode) {
	user := new(BaseUser)
	has, err1 := m.sqlEngine.Where("name=?", form.Name).Get(user)
	if err1 == nil && has {
		md5sum := md5.Sum([]byte(user.Password + m.app.GetSettings().PrivateKey))
		if user.Password == hex.Dump(md5sum[:]) {
			conn := m.redisPool.Get()
			_, err1 := conn.Do("SET", SessionId, user.Id)
			if err1 != nil {
				log.Error("operate redis error")
				return "", base.ErrInternal
			}
			return "success", base.ErrNil
		}
	}
	return "", base.ErrLoginFail
}

func (m *ModuleUser) Regiester(SessionId string, form RegisterForm) (result string, err *base.ErrorCode) {

	if len(form.Name) < 8 || len(form.Password) < 8 {
		return "", base.ErrParamNotAllow
	}
	user := &BaseUser{}
	has, err1 := m.sqlEngine.Where("name=?", form.Name).Get(user)
	if err1 != nil {
		log.Error(err1.Error())
		return "", base.ErrInternal
	}

	if has {
		return "", base.NewError(1000, "Account has been taken")
	}
	md5sum := md5.Sum([]byte(user.Password + m.app.GetSettings().PrivateKey))
	user = &BaseUser{
		Name:      form.Name,
		Password:  hex.Dump(md5sum[:]),
		CreatTime: time.Now().Unix(),
		UserID:    uuid.Rand().Hex(),
	}
	affected, err3 := m.sqlEngine.Insert(user)
	if err3 != nil {
		log.Error(err1.Error())
		return "", base.ErrInternal
	}
	log.Info("%v", affected)
	return "success", base.ErrNil
}

func (m *ModuleUser) GetVerifyCode(SessionId string, form RegisterForm) (result string, err *base.ErrorCode) {
	randString := uuid.RandNumbers(6)
	conn := m.redisPool.Get()
	_, err1 := conn.Do("SET", form.Name, randString)

	if err1 != nil {
		log.Error("operate redis error")
		return "", base.ErrInternal
	}
	_, err2 := conn.Do("EXPIRE", form.Name, 3600)
	if err2 != nil {
		log.Error(err2.Error())
	}

	return randString, base.ErrNil
}
