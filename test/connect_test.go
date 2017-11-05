package test

import (
	"testing"
	"github.com/GodSlave/MyGoServer/testbase"
	"github.com/GodSlave/MyGoServer/base"
	"fmt"
	"time"
)

func TestMultiConnect(t *testing.T) {

	client, response := testbase.InitConnect(t)
	client1, response1 := testbase.InitConnect(t)
	var err error
	user := &base.BaseUser{
		Name:     "zhanglin",
		Password: "woaini1232",
	}

	err = testbase.LoginI(client, user, response)
	if err != nil {
		t.Fatal(err.Error())
	}

	go func() {
		for {
			select {
			case rep1 := <-response:
				fmt.Printf("rep1 %v  | %v  |  %v  ", rep1.State, string(rep1.Result), rep1.Msg)
				fmt.Println()
			}
		}
	}()
	err = client.Publish(testbase.BuildIPublishMessage(client, nil, "User", "GetSelfInfo"), nil)
	err = testbase.LoginI(client1, user, response1)
	go func() {
		for {
			select {
			case rep2 := <-response1:
				fmt.Printf("rep2 %v  | %v  | %v ", rep2.State, string(rep2.Result), rep2.Msg)
				fmt.Println()
			}
		}
	}()

	if err != nil {
		t.Fatal(err.Error())
	}
	err = client1.Publish(testbase.BuildIPublishMessage(client1, nil, "User", "GetSelfInfo"), nil)
	time.Sleep(200)

	err = client.Publish(testbase.BuildIPublishMessage(client, nil, "User", "GetSelfInfo"), nil)
	if err == nil {
		t.Error(" client not kickout ")
	}
	err = client1.Publish(testbase.BuildIPublishMessage(client1, nil, "User", "GetSelfInfo"), nil)

	time.Sleep(2 * time.Second)
}
