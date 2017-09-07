package user

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/utils"
	"github.com/GodSlave/MyGoServer/log"
)

type ModuleUser struct {
	basemodule.BaseModule
	redisPool *redis.Pool
}

type BaseUser struct {
	Name     string
	password string
	Id       string
}

var Module = func() module.Module {
	newGate := new(ModuleUser)

	return newGate
}

func (m *ModuleUser) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	url := settings.Settings["redis"].(string)
	m.redisPool = utils.GetRedisFactory().GetPool(url)

	m.GetServer().RegisterGO("Login", m.Login)
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

func (m *ModuleUser) Login(SessionId string, msg map[string]interface{}) (result string, err string) {
	name := msg["Name"].(string)
	password := msg["Password"].(string)
	if name == password {
		conn := m.redisPool.Get()
		_, err1 := conn.Do("SET", SessionId, name)
		_, err2 := conn.Do("SET", name, SessionId)
		if err1 != nil || err2 != nil {
			log.Error("operate redis error")
		}
	}else {
		return "","password error"
	}

	return "success", ""
}
