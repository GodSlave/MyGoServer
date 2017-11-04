package main

import (
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/module/user"
)

var(
	Version string
	Build   string
)

func main() {
		app := app2.NewApp()
		app.Run(gate.Module(), user.Module())
}

