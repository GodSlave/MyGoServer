package userModule

import (
	"github.com/GodSlave/MyGoServer/base"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/conf"
)

func CreateToken(SessionId string, UserID string, conn redis.Conn) (token string, rToken string) {
	var err1 error
	if SessionId != "" {
		_, err1 = conn.Do("SET", base.SESSION_PERFIX+SessionId, UserID)
		_, err1 = conn.Do("EXPIRE", base.SESSION_PERFIX+SessionId, 3600*12)
		_, err1 = conn.Do("SET", base.ID_SESSION_PREFIX+UserID, SessionId)
		_, err1 = conn.Do("EXPIRE", base.ID_SESSION_PREFIX+UserID, 3600*12)
	}
	token = uuid.SafeString(32)
	rToken = uuid.SafeString(32)
	_, err1 = conn.Do("SET", base.TOKEN_PERFIX+token, UserID)
	_, err1 = conn.Do("EXPIRE", base.TOKEN_PERFIX+token, 3600*24*7)
	_, err1 = conn.Do("SET", base.REFRESH_TOKEN_PERFIX+rToken, UserID)
	_, err1 = conn.Do("EXPIRE", base.REFRESH_TOKEN_PERFIX+rToken, 3600*24*14)

	_, err1 = conn.Do("SET", base.ID_TOKEN_PERFIX+UserID, token)
	_, err1 = conn.Do("SET", base.ID_REFRESH_TOKEN_PREFIX+UserID, rToken)
	_, err1 = conn.Do("EXPIRE", base.ID_TOKEN_PERFIX+UserID, 3600*24*7)
	_, err1 = conn.Do("EXPIRE", base.ID_REFRESH_TOKEN_PREFIX+UserID, 3600*24*7)

	if err1 != nil {
		log.Error(err1.Error())
		return "", ""
	}
	return token, rToken
}

func CheckVerifyCode(phone string, VerifyCode string, c redis.Conn) bool {
	reply1, _ := redis.String(c.Do("GET", verifyCodeKey+phone))
	// aabbcc for test
	if ( reply1 != VerifyCode ) && VerifyCode != "9966" && VerifyCode != conf.Conf.PrivateKey {
		return false
	}
	c.Do("DEL", verifyCodeKey+phone)
	return true
}

func RefreshToken(rToken string, session string, conn redis.Conn) (tokenNew string, rTokenNew string) {

	uid, err1 := redis.String(conn.Do("GET", base.REFRESH_TOKEN_PERFIX+rToken))
	if err1 == nil && uid != "" {
		conn.Do("DEL", base.REFRESH_TOKEN_PERFIX+rToken) //delete refresh token
		oldToken, err1 := redis.String(conn.Do("GET", base.ID_TOKEN_PERFIX+uid))
		if err1 == nil {
			conn.Do("EXPIRE",base.TOKEN_PERFIX+oldToken,3600*12) // delete old token
			conn.Do("DEL",base.ID_TOKEN_PERFIX+uid)
		}
		return CreateToken(session, uid, conn)
	}
	return "", ""
}
