package gate

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/mqtt/sessions"
	"github.com/GodSlave/MyGoServer/mqtt/service"
	"github.com/surgemq/message"
	"encoding/json"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/gogo/protobuf/proto"
	"github.com/GodSlave/MyGoServer/utils/aes"
)

var RPC_PARAM_SESSION_TYPE = "SESSION"

type Gate struct {
	basemodule.BaseModule
	svr *service.Server
}

type MsgFormat struct {
	Module string      `json:"module"`
	Func   string      `json:"func"`
	Params interface{} `json:"params"`
}

var Module = func() module.Module {
	newGate := new(Gate)
	return newGate
}

func (m *Gate) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Gate"
}
func (m *Gate) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

func (m *Gate) OnInit(app module.App, settings *conf.ModuleSettings) {
	//read setting
	m.BaseModule.OnInit(m, app, settings) //这是必须的
}

func (m *Gate) Run(closeSig chan bool) {
	m.svr = &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
		MsgAgent:         m,
		Services:         make(map[string]*service.Service),
	}
	// Listen and serve connections at localhost:1883
	m.svr.ListenAndServe("tcp://0.0.0.0:1883")
	<-closeSig
	m.svr.Close()

}

func (m *Gate) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}

func (m *Gate) robot(session sessions.Session, r []byte) (result string, err string) {
	//time.Sleep(1)
	return string(r), ""
}

func (m *Gate) Process(msg *message.PublishMessage, sess *sessions.Session) bool {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Gate  OnRecover error [%s]", r)
		}
	}()
	topic := msg.Topic()
	log.Info(msg.String())
	if topic[0] == 'i' {
		return m.progressJsonMessage(msg, sess)
	} else if topic[0] == 'f' {
		return m.progressProtoMessage(msg, sess)
	}

	return true
}

func (m *Gate) progressJsonMessage(msg *message.PublishMessage, sess *sessions.Session) bool {
	toResult := func(Topic []byte, Result []byte, error *base.ErrorCode, packetId uint16) (err error) {
		r := &ResultInfo{
			Error:     error.Desc,
			Result:    Result,
			ErrorCode: error.ErrorCode,
		}
		b, err := json.Marshal(r)
		if err == nil {
			err := m.WriteMsg(Topic, b, packetId, sess)
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			r = &ResultInfo{
				Error:  err.Error(),
				Result: nil,
			}
			log.Error(err.Error())

			br, _ := json.Marshal(r)
			err := m.WriteMsg(Topic, br, packetId, sess)
			if err != nil {
				log.Error(err.Error())
			}
		}
		return
	}
	topic := msg.Topic()
	payload := msg.Payload()
	if m.GetApp().GetSettings().Secret {
		aesCipher, _ := aes.NewAesEncrypt(sess.AesKey)
		payload, err := aesCipher.Decrypt(payload)
		log.Info(string(payload))
		if err != nil {
			log.Error(err.Error())
			return false
		}
	}

	var msgContent MsgFormat
	packetId := msg.PacketId()
	err := json.Unmarshal(payload, &msgContent)
	if err == nil {
		arg1, err := json.Marshal(msgContent.Params)
		if err != nil {
			log.Error("msg param format error %s", err.Error())
			toResult(topic, nil, base.ErrParamParseFail, packetId)
		}
		serverSession, err := m.App.GetRouteServers(msgContent.Module, "")
		if err != nil {
			log.Error("Service(type:%s) not found", msgContent.Module)
			toResult(topic, nil, base.ErrNotFound, packetId)
			return false
		}
		result, e := serverSession.CallArgs(msgContent.Func, sess.Id, arg1)
		toResult(topic, result, e, packetId)
	} else {
		log.Error("msg format error %s", err.Error())
		toResult(topic, nil, base.ErrParamParseFail, packetId)
		return false
	}
	return true
}

func (m *Gate) progressProtoMessage(msg *message.PublishMessage, sess *sessions.Session) bool {
	toResult := func(Topic []byte, Result []byte, error *base.ErrorCode, packetId uint16, ) (err error) {
		r := &ResultInfo{
			Error:     error.Desc,
			Result:    Result,
			ErrorCode: error.ErrorCode,
		}
		b, err := proto.Marshal(r)
		if err == nil {
			err := m.WriteMsg(Topic, b, packetId, sess)
			if (err != nil) {
				log.Error(err.Error())
			}
		} else {
			r = &ResultInfo{
				Error:  err.Error(),
				Result: nil,
			}
			log.Error(err.Error())

			br, _ := proto.Marshal(r)
			err := m.WriteMsg(Topic, br, packetId, sess)
			if err != nil {
				log.Error(err.Error())
			}
		}
		return
	}
	topic := msg.Topic()
	payload := msg.Payload()

	if m.GetApp().GetSettings().Secret {
		aesCipher, _ := aes.NewAesEncrypt(sess.AesKey)
		var err error
		payload, err = aesCipher.Decrypt(payload)
		if err != nil {
			log.Error(err.Error())
			return false
		}
	}

	packetId := msg.PacketId()
	serverSession, err := m.App.GetByteRouteServers(payload[0], "")
	if err != nil {
		log.Error("Service(type:%x) not found", payload[0])
		toResult(topic, nil, base.ErrNotFound, packetId)
		return false
	}

	result, e := serverSession.CallByteArgs(payload[1], sess.Id, payload[3:])
	realResult := make([]byte, len(result)+3)
	copy(realResult[0:3], payload[0:3])
	copy(realResult[3:], result)
	toResult(topic, realResult, e, packetId)
	return true
}

func (m *Gate) WriteMsg(topic []byte, body []byte, packetId uint16, sess *sessions.Session) error {
	publish := message.NewPublishMessage()
	publish.SetTopic(topic)
	log.Info(string(topic))
	log.Info(string(body))
	if (m.GetApp().GetSettings().Secret) {
		aesCipher, _ := aes.NewAesEncrypt(sess.AesKey)
		var err error
		body, err = aesCipher.EncryptBytes(body)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	publish.SetPayload(body)
	publish.SetPacketId(packetId)
	publish.SetQoS(1)

	return m.svr.PublishToClient(publish, sess.Id, nil)
}

func (m *Gate) DisConnect(sess *sessions.Session) {
	log.Info("%s disconnect ", sess.Id)
}
