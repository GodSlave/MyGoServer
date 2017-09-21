package user

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/utils"
	"github.com/go-xorm/xorm"
	"crypto/md5"
	"encoding/hex"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/base"
	"time"
	"github.com/GodSlave/MyGoServer/utils/uuid"
)

type ModuleUser struct {
	basemodule.BaseModule
	redisPool *redis.Pool
	sqlEngine *xorm.Engine
	app       module.App
}

var verifyCodeKey = "VerifyCode"


var Module = func() module.Module {
	newGate := new(ModuleUser)

	return newGate
}

func (m *ModuleUser) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	m.redisPool =app.GetRedis()
	m.sqlEngine = app.GetSqlEngine()
	m.GetServer().RegisterGO("Login", 1, m.Login)
	m.GetServer().RegisterGO("Register", 2, m.Register)
	m.GetServer().RegisterGO("GetVerifyCode", 3, m.GetVerifyCode)
	m.app = app

	var user = &base.BaseUser{}
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

func (m *ModuleUser) Login(SessionId string, form *UserLoginRequest) (result *UserLoginResponse, err *base.ErrorCode) {
	user := new(base.BaseUser)
	has, err1 := m.sqlEngine.Where("name=?", form.Username).Get(user)
	if err1 == nil && has {
		md5sum := md5.Sum([]byte(user.Password + m.app.GetSettings().PrivateKey))
		if user.Password == hex.EncodeToString(md5sum[:]) {
			conn := m.redisPool.Get()
			_, err1 := conn.Do("SET", base.SESSION_PERFIX+SessionId, user.Id)
			_, err1 = conn.Do("EXPIRE", base.SESSION_PERFIX+SessionId, 3600*24)

			token := uuid.RandNumbers(32)
			rToken := uuid.RandNumbers(32)

			_, err1 = conn.Do("SET", base.TOKEN_PERFIX+token, user.Id)
			_, err1 = conn.Do("EXPIRE", base.TOKEN_PERFIX+token, 3600*24*7)
			_, err1 = conn.Do("SET", base.REFRESH_TOKEN_PERFIX+rToken, user.Id)
			_, err1 = conn.Do("EXPIRE", base.REFRESH_TOKEN_PERFIX+rToken, 3600*24*14)

			if err1 != nil {
				log.Error("operate redis error")
				return nil, base.ErrInternal
			}

			loginReq := &UserLoginResponse{

				UserTokenData: &UserTokenData{
					Token:        token,
					RefreshToken: rToken,
					ExpireAt:     time.Now().AddDate(0, 0, 7).Unix(),
				},
			}

			return loginReq, base.ErrNil
		}
	}
	return nil, base.ErrLoginFail
}

func (m *ModuleUser) Register(SessionId string, form *UserRegisterRequest) (result *UserRegisterResponse, err *base.ErrorCode) {

	if len(form.Username) < 8 || len(form.Password) < 8 {
		return nil, base.ErrNameOrPwdShort
	}
	user := &base.BaseUser{}
	has, err1 := m.sqlEngine.Where("name=?", form.Username).Get(user)
	if err1 != nil {
		log.Error(err1.Error())
		return nil, base.ErrInternal
	}

	if has {
		return nil, base.ErrAccountBeenTaken
	}
	md5sum := md5.Sum([]byte(user.Password + m.app.GetSettings().PrivateKey))
	user = &base.BaseUser{
		Name:      form.Username,
		Password:  hex.EncodeToString(md5sum[:]),
		CreatTime: time.Now().Unix(),
		UserID:    uuid.Rand().Hex(),
	}
	affected, err3 := m.sqlEngine.Insert(user)
	if err3 != nil {
		log.Error(err1.Error())
		return nil, base.ErrInternal
	}
	log.Info("%v", affected)
	return &UserRegisterResponse{
		Result: "success",
	}, base.ErrNil
}

func (m *ModuleUser) GetVerifyCode(SessionId string, form UserGetVerifyCodeRequest) (result *UserGetVerifyCodeResponse, err *base.ErrorCode) {
	randString := uuid.RandNumbers(6)
	conn := m.redisPool.Get()
	var err1 error
	_, err1 = conn.Do("SET", verifyCodeKey+form.PhoneNumber, randString)
	_, err1 = conn.Do("EXPIRE", verifyCodeKey+form.PhoneNumber, 3600)

	if err1 != nil {
		log.Error("operate redis error")
		return nil, base.ErrInternal
	}

	return nil, base.ErrNil
}

func (m *ModuleUser) GetSelfInfo(user *base.BaseUser, form UserGetSelfInfoRequest) (result *UserGetSelfInfoResponse, err *base.ErrorCode) {
	if user == nil {
		err = base.ErrNeedLogin
		return
	}

	result = &UserGetSelfInfoResponse{
		UserData: &UserData{
			UserName: user.Name,
			UserID:   user.UserID,
		},
	}

	return
}
