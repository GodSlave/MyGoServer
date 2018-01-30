package main

import (
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/userModule"
)

func main() {
	app := app2.NewApp()
	app.Run(userModule.Module())
}
