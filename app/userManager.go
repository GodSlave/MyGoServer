package app

import (
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/utils/lru"
)

type DefaultUserManager struct {
	module.UserManager
	app              module.App
	users            *lru.Cache
	loginCallBack    module.UserEventCallBack
	logOutCallBack   module.UserEventCallBack
	registerCallBack module.UserEventCallBack
}

func InitUserManager(app module.App, lruSize int32) (um module.UserManager) {
	duf := &DefaultUserManager{}
	duf.app = app
	duf.users = lru.New(int(lruSize))
	return duf
}

func (um *DefaultUserManager) OnUserLogOut(user *base.BaseUser) {
	if um.logOutCallBack != nil {
		um.logOutCallBack(user)
	}
}
func (um *DefaultUserManager) OnUserLogin(user *base.BaseUser) {
	if um.loginCallBack != nil {
		um.loginCallBack(user)
	}
}

func (um *DefaultUserManager) OnUserRegister(user *base.BaseUser) {
	if um.registerCallBack != nil {
		um.registerCallBack(user)
	}
}

func (um *DefaultUserManager) CheckUserLogin(sessionId string) {
	um.VerifyUser(sessionId)

}
func (um *DefaultUserManager) CheckUserLogOut(sessionId string) {
	_, exit := um.users.Get(sessionId)
	if exit {
		um.users.Remove(sessionId)
	}
}
func (um *DefaultUserManager) VerifyUser(sessionId string) (user *base.BaseUser) {
	obj, exit := um.users.Get(sessionId)
	user = obj.(*base.BaseUser)
	if exit {
		return
	}
	uid := um.VerifyUserID(sessionId)
	if uid == "" {
		return nil
	}
	user = &base.BaseUser{
		UserID: uid,
	}
	result, err := um.app.GetSqlEngine().Get(user)
	if err != nil {
		return nil
	}
	if !result {
		return nil
	}
	um.users.Add(sessionId, user)
	return
}
func (um *DefaultUserManager) VerifyUserID(sessionId string) (userID string) {
	c := um.app.GetRedis().Get()
	has, err := redis.Bool(c.Do("EXISTS", base.TOKEN_PERFIX+sessionId))
	if err == nil && has {
		userID, _ = redis.String(c.Do("GET", base.TOKEN_PERFIX+sessionId))
		return
	}

	has, err = redis.Bool(c.Do("EXISTS", base.SESSION_PERFIX+sessionId))
	if err == nil && has {
		userID, _ = redis.String(c.Do("GET", base.SESSION_PERFIX+sessionId))
		return
	}
	return ""
}

func (um *DefaultUserManager) SetLoginCallBack(callback module.UserEventCallBack) {
	um.loginCallBack = callback
}
func (um *DefaultUserManager) SetRegisterCallBack(callback module.UserEventCallBack) {
	um.registerCallBack = callback
}
func (um *DefaultUserManager) SetLogoutCallBack(callback module.UserEventCallBack) {
	um.registerCallBack = callback
}
