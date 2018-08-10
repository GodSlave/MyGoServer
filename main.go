package main

import (
	"github.com/GodSlave/MyGoServer/template"
	"github.com/GodSlave/MyGoServer/base"
)

var (
	Version string
	Build   string
)

func main() {
	//app := app2.NewApp()
	//web_Module := web.Module()
	//webModule := web_Module.(*web.ModuleWeb)
	//webModule.Router = gin.Default()
	//app.Run(gate.Module(), userModule.Module(), httpGate.Module(), web_Module)
	template.BuildModel(base.BaseUser{},"User")
}
