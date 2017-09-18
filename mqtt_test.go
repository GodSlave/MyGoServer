package main

import (
	"github.com/GodSlave/MyGoServer/mqtt/service"
	"testing"
	app2 "github.com/GodSlave/MyGoServer/app"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/module/user"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"fmt"
	"github.com/surgemq/message"
	"github.com/GodSlave/MyGoServer/utils/aes"
	"crypto/md5"
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

func TestProtoLogin(t *testing.T) {
	// Instantiates a new Client
	c := &service.Client{}

	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetCleanSession(true)
	msg.SetVersion(0x4)
	id := uuid.SafeString(33);
	fmt.Println(id)
	msg.SetClientId([]byte(id))
	msg.SetKeepAlive(10)

	// Connects to the remote server at 127.0.0.1 port 1883
	err := c.Connect("tcp://127.0.0.1:1883", msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Creates a new SUBSCRIBE message to subscribe to topic "abc"
	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte("s"), 1)

	// Subscribes to the topic by sending the message. The first nil in the function
	// call is a OnCompleteFunc that should handle the SUBACK message from the server.
	// Nil means we are ignoring the SUBACK messages. The second nil should be a
	// OnPublishFunc that handles any messages send to the client because of this
	// subscription. Nil means we are ignoring any PUBLISH messages for this topic.

	chanOut := make(chan int, 3)
	c.Subscribe(submsg, func(msg, ack message.Message, err error) error {
		fmt.Println(ack.Name())
		return nil
	}, func(msg *message.PublishMessage) error {
		fmt.Println(msg.String())
		if string(msg.Topic()) == ("s") {

			pubmsg := message.NewPublishMessage()
			pubmsg.SetTopic([]byte("i"))
			body := `{"module":"User","func":"Login","params":{"name":"zhang","password":"123456"}}`
			datas := []byte(body)
			aes, _ := aes.NewAesEncrypt(md5result[:])
			fmt.Println("%v", datas)
			datas, _ = aes.EncryptBytes(datas)
			fmt.Println("%v", datas)
			ddatas, _ := aes.Decrypt(datas)
			fmt.Println("%v", ddatas)
			fmt.Println(string(ddatas))
			pubmsg.SetPayload(datas)
			pubmsg.SetQoS(1)
			c.Publish(pubmsg, nil)

		}
		chanOut <- 1
		return nil
	})

	// Publishes to the server by sending the message

	fmt.Println("wait msg")
	<-chanOut
	<-chanOut
	<-chanOut
	fmt.Println("end")
	//time.Sleep(1000)
	//c.Disconnect()
}
