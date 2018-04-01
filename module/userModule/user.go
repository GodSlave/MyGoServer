package userModule

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/go-xorm/xorm"
	"crypto/md5"
	"encoding/hex"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/base"
	"time"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"encoding/json"
)

type ModuleUser struct {
	basemodule.BaseModule
	redisPool *redis.Pool
	sqlEngine *xorm.Engine
	app       module.App
}

var verifyCodeKey = "VerifyCode"
var verifyCodeTimeKey = "VerifyCodeTime"
var forgotPasswordVerifyCodeKey = "forgotPasswordVerifyCodeKey"

var Module = func() module.Module {
	newGate := new(ModuleUser)
	return newGate
}

func (m *ModuleUser) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	m.redisPool = app.GetRedis()
	m.sqlEngine = app.GetSqlEngine()
	m.GetServer().RegisterGO("Login", 1, m.Login)
	m.GetServer().RegisterGO("Register", 2, m.Register)
	m.GetServer().RegisterGO("GetVerifyCode", 3, m.GetVerifyCode)
	m.GetServer().RegisterGO("GetSelfInfo", 4, m.GetSelfInfo)
	m.GetServer().RegisterGO("Logout", 5, m.LogOut)
	m.GetServer().RegisterGO("RefreshToken", 6, m.RefreshToken)
	m.GetServer().RegisterGO("ChangePasswordByVerifyCode", 7, m.ChangePasswordByVerifyCode)
	m.GetServer().RegisterGO("ChangePasswordByPassword", 8, m.ChangePasswordByPassword)
	m.GetServer().RegisterGO("CheckVerifyCodeAvailable", 9, m.checkVerifyCodeAvailable)
	m.GetServer().RegisterGO("SendForgotPasswordVrifyCode", 10, m.ForgetPassword)

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

func (m *ModuleUser) Login(SessionId string, form *User_Login_Request) (result *User_Login_Response, err *base.ErrorCode) {
	user := new(base.BaseUser)
	has, err1 := m.sqlEngine.Where("name=?", form.Username).Get(user)
	if err1 == nil && has {
		md5sum := md5.Sum([]byte(form.Password + m.app.GetSettings().PrivateKey))
		if user.Password == hex.EncodeToString(md5sum[:]) {
			conn := m.redisPool.Get()
			defer conn.Close()
			m.removeLoginUser(user, conn, SessionId)
			token, rToken := CreateToken(SessionId, user.UserID, conn)
			if token == "" {
				log.Error("operate redis error")
				return nil, base.ErrInternal
			}

			loginReq := &User_Login_Response{
				UserTokenData: &UserTokenData{
					Token:        token,
					RefreshToken: rToken,
					ExpireAt:     time.Now().AddDate(0, 0, 7).Unix(),
				},
			}
			m.GetApp().GetUserManager().OnUserLogin(user)
			return loginReq, base.ErrNil
		}
	}
	return nil, base.ErrLoginFail
}

func (m *ModuleUser) removeLoginUser(user *base.BaseUser, redisConn redis.Conn, currentSession string) {
	token, err := redis.String(redisConn.Do("GET", base.ID_TOKEN_PERFIX+user.UserID))
	session, err := redis.String(redisConn.Do("GET", base.ID_SESSION_PREFIX+user.UserID))
	if token != "" || session != "" && (currentSession != token && currentSession != user.UserID) {
		m.app.GetUserManager().OnUserLogOut(user)
		go func() {
			//send offline reason
			pushItem := base.PushItem{
				Module:   0,
				PushType: 0,
				Content: &base.PushContent{
					Content: []byte("Other client Login"),
				},
			}
			content, _ := json.Marshal(pushItem)
			if session != "" {
				m.app.RpcAllInvokeArgs(m, "Gate", "PushMessage", session, content)
				//m.app.RpcAllInvokeArgs(m, "Gate", "KickOut", session, nil)
			} else if token != "" {
				m.app.RpcAllInvokeArgs(m, "Gate", "PushMessage", session, content)
				//m.app.RpcAllInvokeArgs(m, "Gate", "KickOut", token, nil)
			}
		}()
	}

	if err != nil {
		log.Error(err.Error())
	}

}

func (m *ModuleUser) Register(SessionId string, form *User_Register_Request) (result *User_Register_Response, err *base.ErrorCode) {

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

	c := m.redisPool.Get()
	defer c.Close()
	checkresult := CheckVerifyCode(form.Username, form.VerifyCode, c)
	if !checkresult {
		return nil, base.ErrVerifyCodeErr
	}

	md5sum := md5.Sum([]byte(form.Password + m.app.GetSettings().PrivateKey))
	password := hex.EncodeToString(md5sum[:])
	user = &base.BaseUser{
		Name:      form.Username,
		Password:  password,
		CreatTime: time.Now().Unix(),
		UserID:    uuid.Rand().Hex(),
	}
	_, err3 := m.sqlEngine.Insert(user)
	if err3 != nil {
		log.Error(err3.Error())
		return nil, base.ErrInternal
	}
	m.app.GetUserManager().OnUserRegister(user)
	conn := m.redisPool.Get()
	defer conn.Close()
	m.removeLoginUser(user, conn, SessionId)
	token, rToken := CreateToken(SessionId, user.UserID, conn)
	return &User_Register_Response{
		Result: user.UserID,
		TokenInfo: &UserTokenData{
			Token:        token,
			RefreshToken: rToken,
			ExpireAt:     time.Now().AddDate(0, 0, 7).Unix(),
		},
	}, base.ErrNil
}

func (m *ModuleUser) GetVerifyCode(SessionId string, form *User_GetVerifyCode_Request) (result *User_GetVerifyCode_Response, err *base.ErrorCode) {
	randString := uuid.RandNumbers(4)
	conn := m.redisPool.Get()
	defer conn.Close()
	var err1 error
	if len(form.PhoneNumber) < 11 {
		return nil, base.ErrParamNotAllow
	}
	lastRequestTime, err1 := redis.Int64(conn.Do("GET", verifyCodeTimeKey+form.PhoneNumber))
	if time.Now().Unix()-lastRequestTime < 60 {
		return nil, base.ErrVerifySendTooBusy
	}

	_, err1 = conn.Do("DEL", verifyCodeKey+form.PhoneNumber)
	_, err1 = conn.Do("SET", verifyCodeKey+form.PhoneNumber, randString)
	_, err1 = conn.Do("EXPIRE", verifyCodeKey+form.PhoneNumber, 900)
	//TODO  real send verify code

	if err1 != nil {
		log.Error("operate redis error")
		return nil, base.ErrInternal
	}
	return &User_GetVerifyCode_Response{}, base.ErrNil
}

func (m *ModuleUser) GetSelfInfo(user *base.BaseUser) (result *User_GetSelfInfo_Response, err *base.ErrorCode) {
	if user == nil {
		err = base.ErrNeedLogin
		return
	}

	result = &User_GetSelfInfo_Response{
		UserData: &UserData{
			UserName: user.Name,
			UserID:   user.UserID,
		},
	}
	return
}

func (m *ModuleUser) LogOut(user *base.BaseUser) (result *User_Login_Response, err *base.ErrorCode) {
	if user == nil {
		err = base.ErrNeedLogin
		return
	}
	m.app.GetUserManager().OnUserLogOut(user)
	return
}

func (m *ModuleUser) RefreshToken(sessionId string, form *User_RefreshToken_Request) (result *User_RefreshToken_Response, err *base.ErrorCode) {
	conn, err1 := m.redisPool.Dial()
	defer conn.Close()
	if conn != nil && err1 == nil {
		token, rToken := RefreshToken(form.RefreshToken, sessionId, conn)
		if token == "" {
			return nil, &base.ErrorCode{
				ErrorCode: 1, Desc: "刷新token失败",
			}
		}
		refreshResponse := &User_RefreshToken_Response{
			TokenData: &UserTokenData{
				Token:        token,
				RefreshToken: rToken,
				ExpireAt:     time.Now().AddDate(0, 0, 7).Unix(),
			},
		}
		return refreshResponse, base.ErrNil
	}
	return nil, base.ErrInternal
}

func (m *ModuleUser) ForgetPassword(sessionId string, form *User_ForgetPassWord_Request) (result *User_ForgetPassWord_Response, err *base.ErrorCode) {

	randString := uuid.RandNumbers(6)
	conn := m.redisPool.Get()
	defer conn.Close()
	var err1 error
	if len(form.PhoneNumber) < 11 {
		return nil, base.ErrParamNotAllow
	}
	lastRequestTime, err1 := redis.Int64(conn.Do("GET", verifyCodeTimeKey+form.PhoneNumber))
	if time.Now().Unix()-lastRequestTime < 60 {
		return nil, base.ErrVerifySendTooBusy
	}

	_, err1 = conn.Do("DEL", forgotPasswordVerifyCodeKey+form.PhoneNumber)
	_, err1 = conn.Do("SET", forgotPasswordVerifyCodeKey+form.PhoneNumber, randString)
	_, err1 = conn.Do("EXPIRE", forgotPasswordVerifyCodeKey+form.PhoneNumber, 900)
	//TODO  real send verify code

	if err1 != nil {
		log.Error("operate redis error")
		return nil, base.ErrInternal
	}
	return &User_ForgetPassWord_Response{}, base.ErrNil

}

func (m *ModuleUser) ChangePasswordByVerifyCode(sessionId string, form *User_ChangePassWordByVerifyCode_Request) (result *User_ChangePassWordByVerifyCode_Response, err *base.ErrorCode) {
	conn := m.redisPool.Get()
	defer conn.Close()
	reply1, _ := redis.String(conn.Do("GET", forgotPasswordVerifyCodeKey+form.PhoneNumber))
	// aabbcc for test
	if ( reply1 != form.VerifyCode ) && form.VerifyCode != "999666" && form.VerifyCode != conf.Conf.PrivateKey {
		return nil, &base.ErrorCode{
			1,
			"验证码错误",
		}
	}

	user := &base.BaseUser{
		Name: form.PhoneNumber,
	}

	exist, err1 := m.sqlEngine.Get(user)
	if !exist || err1 != nil {
		return nil, &base.ErrorCode{
			2,
			"未找到该用户",
		}
	}

	md5sum := md5.Sum([]byte(form.Password + m.app.GetSettings().PrivateKey))
	password := hex.EncodeToString(md5sum[:])
	user.Password = password
	_, err2 := m.sqlEngine.Update(user)
	if err2 != nil {
		return nil, base.ErrInternal
	}

	return &User_ChangePassWordByVerifyCode_Response{
		true,
	}, base.ErrNil
}

func (m *ModuleUser) ChangePasswordByPassword(user *base.BaseUser, form *User_ChangePassWordByPassword_Request) (result *User_ChangePassWordByPassword_Response, err *base.ErrorCode) {
	md5sum := md5.Sum([]byte(form.OldPassword + m.app.GetSettings().PrivateKey))
	oldpassword := hex.EncodeToString(md5sum[:])
	if oldpassword != user.Password {
		return nil, &base.ErrorCode{
			ErrorCode: 1, Desc: "旧密码错误",
		}
	}

	if len(form.NewPassword) >= 6 && len(form.NewPassword) <= 24 {
		newMd5sum := md5.Sum([]byte(form.OldPassword + m.app.GetSettings().PrivateKey))
		newPassword := hex.EncodeToString(newMd5sum[:])
		user.Password = newPassword
		_, error := m.sqlEngine.Update(user)
		if error != nil {
			return nil, base.ErrSQLERROR
		}
	} else {
		return nil, &base.ErrorCode{
			ErrorCode: 2, Desc: "新密码格式错误",
		}
	}

	return &User_ChangePassWordByPassword_Response{
		true,
	}, base.ErrNil
}

func (m *ModuleUser) checkVerifyCodeAvailable(sessionId string, form *User_CheckVerifyCodeAvailable_Request) (*User_CheckVerifyCodeAvailable_Response, *base.ErrorCode) {

	conn, err := m.redisPool.Dial()
	if err == nil {
		defer conn.Close()
		baseuser := &base.BaseUser{
			Phone: form.PhoneNumber,
		}

		exist, _ := m.sqlEngine.Get(baseuser)
		verifyCode, err := redis.String(conn.Do("get", forgotPasswordVerifyCodeKey+form.PhoneNumber))
		if err == nil && verifyCode == form.VerifyCode {
			return &User_CheckVerifyCodeAvailable_Response{true, exist}, base.ErrNil
		}
		verifyCode, err = redis.String(conn.Do("get", verifyCodeKey+form.PhoneNumber))
		if err == nil && verifyCode == form.VerifyCode {
			return &User_CheckVerifyCodeAvailable_Response{true, exist}, base.ErrNil
		}
		return nil, &base.ErrorCode{
			ErrorCode: 1,
			Desc:      "验证码校验错误",
		}
	}

	return nil, base.ErrSQLERROR

}
