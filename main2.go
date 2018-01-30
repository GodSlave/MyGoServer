package main

import (
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
)

func main() {
	app := app2.NewApp()
	app.Run(gate.Module())
}
