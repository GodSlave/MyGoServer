package main

import (
	"github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/module/httpGate"
	"github.com/GodSlave/MyGoServer/module/userModule"
	"github.com/GodSlave/MyGoServer/module/web"
	"github.com/gin-gonic/gin"
)

var (
	Version string
	Build   string
)

func main() {
	app := app.NewApp()
	web_Module := web.Module()
	webModule := web_Module.(*web.ModuleWeb)
	webModule.Router = gin.Default()
	app.Run(gate.Module(), userModule.Module(), httpGate.Module(), web_Module)
	//template.BuildModel(base.BaseUser{},"User")
}
