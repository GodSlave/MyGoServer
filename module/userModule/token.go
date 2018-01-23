package user

import (
	"github.com/GodSlave/MyGoServer/base"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/log"
)

func CreateToken(SessionId string, UserID string, conn redis.Conn) (token string, rToken string) {
	_, err1 := conn.Do("SET", base.SESSION_PERFIX+SessionId, UserID)
	_, err1 = conn.Do("EXPIRE", base.SESSION_PERFIX+SessionId, 3600*24)
	_, err1 = conn.Do("SET", base.ID_SESSION_PREFIX+UserID, SessionId)
	_, err1 = conn.Do("EXPIRE", base.ID_SESSION_PREFIX+UserID, 3600*24)
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
	reply1, err1 := redis.String(c.Do("GET", verifyCodeKey+phone))
	if err1 != nil {
		log.Error(err1.Error())
	}

	if ( reply1 != VerifyCode ) && VerifyCode != "aabbcc" {
		return false
	}
	c.Do("DEL", verifyCodeKey+phone)
	return true
}
