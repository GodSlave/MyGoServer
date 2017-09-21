package appModule

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/base"
)

type AppModule struct {
	basemodule.BaseModule
}

var Module = func() module.Module {
	newApp := new(AppModule)
	return newApp
}

func (m *AppModule) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "AppModule"
}

func (m *AppModule) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

func (m *AppModule) OnInit(app module.App, settings *conf.ModuleSettings) {
	//read setting
	m.BaseModule.OnInit(m, app, settings) //这是必须的
	m.GetServer().Register("OnUserLogOut", 1, m.OnUserLogOut)
}

func (m *AppModule) Run(closeSig chan bool) {
	<-closeSig
}

func (m *AppModule) OnUserLogOut(sessionID string) (string, *base.ErrorCode) {
	m.App.OnUserLogOut(sessionID)
	return "", base.ErrNil
}
