package test

import (
	"github.com/GodSlave/MyGoServer/mqtt/service"
	"testing"
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/module/user"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"fmt"
	"github.com/GodSlave/MyGoServer/base"
	"time"
	"github.com/GodSlave/MyGoServer/testbase"
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
	go svr.ListenAndServe("tcp://0.0.0.0:1884")
}

func TestRun(t *testing.T) {
	app := app2.NewApp()
	app.Run(gate.Module(), user.Module())
}

func TestUUID(t *testing.T) {
	fmt.Println(uuid.SafeString(20))
	fmt.Println(uuid.RandNumbers(20))
}

func TestJson(t *testing.T) {

	c := testbase.InitClient()
	time.Sleep(1 * time.Second)
	checkChan := make(chan *gate.AllResponse)
	testbase.SubI(c, checkChan)
	time.Sleep(1 * time.Second)
	user := &base.BaseUser{
		Name:     "zhanglin",
		Password: "woaini1232",
	}
	//RegisterI(c, user, checkChan)
	testbase.LoginI(c, user, checkChan)
	GetSelfInfoI(c, checkChan)

}

func RegisterI(client *service.Client, user1 *base.BaseUser, callback chan *gate.AllResponse) (err error) {
	fmt.Println("start register")
	login := &user.UserRegisterRequest{
		Username:   user1.Name,
		Password:   user1.Password,
		VerifyCode: "aabbcc",
	}
	err = client.Publish(testbase.BuildIPublishMessage(client, login, "User", "Register"), nil)
	var allrespon *gate.AllResponse
	allrespon = <-callback
	fmt.Println("register Response", allrespon.State)
	return err
}

func GetSelfInfoI(client *service.Client, callback chan *gate.AllResponse) (err error) {
	fmt.Println("start register")
	err = client.Publish(testbase.BuildIPublishMessage(client, nil, "User", "GetSelfInfo"), nil)
	var allrespon *gate.AllResponse
	allrespon = <-callback
	fmt.Println("GetSelfInfoI Response", allrespon.State)
	return err
}
