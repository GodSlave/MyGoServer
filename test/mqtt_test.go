package test

import (
	"github.com/GodSlave/MyGoServer/mqtt/service"
	"strconv"
	"testing"
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/module/userModule"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"fmt"
	"github.com/GodSlave/MyGoServer/base"
	"time"
	"github.com/GodSlave/MyGoServer/testbase"
	"sync"
	"github.com/GodSlave/MyGoServer/log"
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
	go svr.ListenAndServe("tcp://0.0.0.0:1883")
}

func TestRun(t *testing.T) {
	app := app2.NewApp()
	app.Run(gate.Module(), userModule.Module())
}

func TestUUID(t *testing.T) {
	fmt.Println(uuid.SafeString(20))
	fmt.Println(uuid.RandNumbers(20))
}


func TestRegister(t *testing.T) {
	startTime :=time.Now().UnixNano()

	numberOfClient := 500
	wg := sync.WaitGroup{}
	wg.Add(numberOfClient)
	for i := 0; i < numberOfClient; i++ {
		go func(index int, wg *sync.WaitGroup) {
			c := testbase.InitClient()
			time.Sleep(1 * time.Second)
			checkChan := make(chan *gate.AllResponse)
			testbase.SubI(c, checkChan)
			time.Sleep(1 * time.Second)
			user := &base.BaseUser{
				Name:     "zhanglin" + strconv.Itoa(index) + uuid.SafeString(5),
				//Password: "woaini1232" + strconv.Itoa(index),
				//Name:     "zhanglin",
				Password: "woaini1232",
			}
			err := RegisterI(c, user, checkChan)
			if err != nil {
				log.Error(err.Error())
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()

	useTime := time.Now().UnixNano()-startTime
	fmt.Println()
	fmt.Println("Use Time ",useTime)
}

func TestGetInfo(t *testing.T) {


	numberOfClient := 5000
	wg := sync.WaitGroup{}
	wg.Add(numberOfClient)
	c := testbase.InitClient()
	time.Sleep(1 * time.Second)
	checkChan := make(chan *gate.AllResponse)
	testbase.SubI(c, checkChan)
	time.Sleep(1 * time.Second)
	user := &base.BaseUser{
		//Name:     "zhanglin" + strconv.Itoa(index) + uuid.SafeString(5),
		//Password: "woaini1232" + strconv.Itoa(index),
		Name:     "zhanglin",
		Password: "woaini1232",
	}
	time.Sleep(1 * time.Second)
	startTime :=time.Now().UnixNano()
	testbase.LoginI(c,user,checkChan)
	for i := 0; i < numberOfClient; i++ {
		go func(index int, wg *sync.WaitGroup) {
			err := GetSelfInfoI(c, checkChan)
			if err != nil {
				log.Error(err.Error())
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()

	useTime := time.Now().UnixNano()-startTime
	fmt.Println()
	fmt.Println("Use Time ",useTime)
	//testbase.LoginI(c, userModule, checkChan)
	//GetSelfInfoI(c, checkChan)
}

func RegisterI(client *service.Client, user1 *base.BaseUser, callback chan *gate.AllResponse) (err error) {
	fmt.Println("start register")
	login := &userModule.User_Register_Request{
		Username:   user1.Name,
		Password:   user1.Password,
		VerifyCode: "9966",
	}
	err = client.Publish(testbase.BuildIPublishMessage(client, login, "User", "Register"), nil)
	log.Info("wait response ")
	var allrespon *gate.AllResponse
	allrespon = <-callback
	fmt.Println("register Response", allrespon.State)
	fmt.Printf("%v", allrespon)
	return err
}

func GetSelfInfoI(client *service.Client, callback chan *gate.AllResponse) (err error) {
	fmt.Println("start getSelfInfo")
	err = client.Publish(testbase.BuildIPublishMessage(client, nil, "User", "GetSelfInfo"), nil)
	var allrespon *gate.AllResponse
	if err == nil {
		allrespon = <-callback
		fmt.Println("GetSelfInfoI Response", allrespon.State)
	} else {
		fmt.Println(err.Error())
	}

	return err
}
