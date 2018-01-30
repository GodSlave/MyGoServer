package app

import (
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/utils/lru"
	"github.com/GodSlave/MyGoServer/log"
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
	conn, err := um.app.GetRedis().Dial()
	if conn != nil && err == nil {
		var token string
		var rToken string
		var err1 error
		token, err1 = redis.String(conn.Do("GET", base.ID_TOKEN_PERFIX+user.UserID))
		rToken, err1 = redis.String(conn.Do("GET", base.ID_REFRESH_TOKEN_PREFIX+user.UserID))
		_, err1 = conn.Do("DEL", base.TOKEN_PERFIX+token)
		_, err1 = conn.Do("DEL", base.REFRESH_TOKEN_PERFIX+rToken)
		_, err1 = conn.Do("DEL", base.ID_TOKEN_PERFIX+user.UserID)
		_, err1 = conn.Do("DEL", base.ID_REFRESH_TOKEN_PREFIX+user.UserID)


		var session string
		session, err1 = redis.String(conn.Do("GET", base.ID_SESSION_PREFIX+user.UserID))
		_, err1 = conn.Do("DEL", base.SESSION_PERFIX+session)
		_, err1 = conn.Do("DEL", base.ID_SESSION_PREFIX+user.UserID)
		if err1 != nil {
			log.Error(err1.Error())
		}

	}
	if um.logOutCallBack != nil {
		um.logOutCallBack(user, um.app)
	}
}
func (um *DefaultUserManager) OnUserLogin(user *base.BaseUser) {
	if um.loginCallBack != nil {
		um.loginCallBack(user, um.app)
	}
}

func (um *DefaultUserManager) OnUserRegister(user *base.BaseUser) {
	log.Info("on userModule register :" + user.Name)
	if um.registerCallBack != nil {
		um.registerCallBack(user, um.app)
	}
}

func (um *DefaultUserManager) OnUserConnect(sessionId string) {
	um.VerifyUser(sessionId)

}
func (um *DefaultUserManager) OnUserDisconnect(sessionId string) {
	user, exit := um.users.Get(sessionId)
	if exit {
		user1 := user.(*base.BaseUser)
		um.users.Remove(sessionId)
		conn, err1 := um.app.GetRedis().Dial()
		if conn != nil && err1 == nil {
			exits, _ := redis.Bool(conn.Do("EXITS", base.SESSION_PERFIX+sessionId))
			if exits {
				_, err1 = conn.Do("DEL", base.SESSION_PERFIX+sessionId)
				_, err1 = conn.Do("DEL", base.ID_SESSION_PREFIX+user1.UserID)
			}

			if err1 != nil {
				log.Error(err1.Error())
			}
		}
	}
}
func (um *DefaultUserManager) VerifyUser(sessionId string) (user *base.BaseUser) {
	obj, exit := um.users.Get(sessionId)
	if exit {
		user = obj.(*base.BaseUser)
		return
	}

	uid := um.VerifyUserID(sessionId)
	if uid == "" {
		log.Info("uid is nil : %s", sessionId)
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
	um.logOutCallBack = callback
}
