package web

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/go-xorm/xorm"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/gin-gonic/gin"
	"github.com/GodSlave/MyGoServer/module"
)

type WebHandlers interface {
	Init(router *gin.Engine)
	SetEngine(sqlEngine *xorm.Engine)
}

type ModuleWeb struct {
	basemodule.BaseModule
	sqlEngine  *xorm.Engine
	app        module.App
	Router     *gin.Engine
	WebModules []WebHandlers
}

var Module = func() module.Module {
	roleModule := new(ModuleWeb)
	return roleModule
}

func (m *ModuleWeb) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.app = app
	m.BaseModule.OnInit(m, app, settings)
	m.sqlEngine = app.GetSqlEngine()
}

func (m *ModuleWeb) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Web"
}

func (m *ModuleWeb) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

func (m *ModuleWeb) Run(closeSig chan bool) {
	m.Router.Static("/public" +
		"", "./public")
	for _, value := range m.WebModules {
		value.Init(m.Router)
		value.SetEngine(m.sqlEngine)
	}
	m.Router.Run(":8090")
	<-closeSig
}

func (m *ModuleWeb) OnDestroy() {


}
