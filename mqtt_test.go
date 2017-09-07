package main

import (
	"github.com/GodSlave/MyGoServer/mqtt/service"
	"testing"
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/module/user"
)

func TestMqtt(t *testing.T) {
	svr := &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
	}
	// Listen and serve connections at localhost:1883
	svr.ListenAndServe("tcp://0.0.0.0:1883")
}

func TestRun(t *testing.T) {
	app := app2.NewApp()
	app.Run(gate.Module(), user.Module())
}
