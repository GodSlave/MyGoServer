package testbase

import (
	"github.com/GodSlave/MyGoServer/mqtt/service"
	"github.com/surgemq/message"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"fmt"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/GodSlave/MyGoServer/module/gate"
	"github.com/GodSlave/MyGoServer/utils/aes"
	"encoding/json"
	"github.com/GodSlave/MyGoServer/module/user"
	"time"
	"strconv"
	"testing"
)

func InitClient() *service.Client {
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
		return nil
	}
	// Creates a new SUBSCRIBE message to subscribe to topic "abc"
	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte("s"), 1)
	chanout := make(chan int, 1)
	c.Subscribe(submsg, func(msg, ack message.Message, err error) error {
		fmt.Println(ack.Name())
		return nil
	}, func(msg *message.PublishMessage) error {
		fmt.Println(msg.String())
		if string(msg.Topic()) == ("s") {
			md5result := base.GetMd5T([]byte(id), msg.Payload())
			c.GetSession().AesKey = md5result
		}
		chanout <- 1
		return nil
	})
	<-chanout
	return c
}

func SubI(client *service.Client, checkChan chan *gate.AllResponse) error {
	fmt.Println("start sub i")
	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte("i"), 1)
	chanout := make(chan int, 1)
	client.Subscribe(submsg, func(msg, ack message.Message, err error) error {
		fmt.Println("ack ", ack)
		chanout <- 1
		return nil
	}, func(msg *message.PublishMessage) error {
		fmt.Println("get response " + msg.String())
		payload := msg.Payload()
		aesCipher, _ := aes.NewAesEncrypt(client.GetSession().AesKey)
		var err error
		payload, err = aesCipher.Decrypt(payload)
		var allRepon = &gate.AllResponse{}
		err = json.Unmarshal(payload, allRepon)
		if err == nil {
			checkChan <- allRepon
		} else {
			fmt.Println(err.Error())
		}
		return err
	})
	<-chanout
	return nil
}

func LoginI(client *service.Client, user1 *base.BaseUser, callback chan *gate.AllResponse) (err error) {
	fmt.Println("start login")
	login := &user.UserLoginRequest{
		Username: user1.Name,
		Password: user1.Password,
	}
	err = client.Publish(BuildIPublishMessage(client, login, "User", "Login"), nil)
	var allrespon *gate.AllResponse
	allrespon = <-callback
	fmt.Println("login Response", allrespon.State)
	return err
}
func BuildIPublishMessage(c *service.Client, value interface{}, module string, method string) *message.PublishMessage {
	msg := gate.MsgFormat{
		Module: module,
		Func:   method,
		Params: value,
	}
	var err error
	data, err := json.Marshal(msg)
	aesCipher, _ := aes.NewAesEncrypt(c.GetSession().AesKey)
	data, err = aesCipher.EncryptBytes(data)
	if err == nil {
		pub := message.NewPublishMessage()
		pub.SetPayload(data)
		pub.SetQoS(1)
		pub.SetTopic([]byte("i"))
		return pub
	}

	return nil
}


func Login(t *testing.T) (*service.Client, *base.BaseUser, chan *gate.AllResponse) {
	var err error
	c := InitClient()
	if c == nil {
		t.Fatal("clint build fail")
	}
	time.Sleep(1 * time.Second)
	checkChan := make(chan *gate.AllResponse)
	SubI(c, checkChan)
	time.Sleep(1 * time.Second)
	user := &base.BaseUser{
		Name:     "zhanglin",
		Password: "woaini1232",
	}
	err = LoginI(c, user, checkChan)
	if err != nil {
		t.Fatal(err.Error())
	}
	return c, user, checkChan
}

func SendMessage(c *service.Client, checkChan chan *gate.AllResponse, publishMessage *message.PublishMessage) (err error) {
	err = c.Publish(publishMessage, nil)
	var allrespon *gate.AllResponse
	allrespon = <-checkChan
	fmt.Println("Get Response", strconv.Itoa(int(allrespon.State))+"  "+allrespon.Msg)
	return err
}