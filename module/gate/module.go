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
)

var RPC_PARAM_SESSION_TYPE = "SESSION"

type Gate struct {
	basemodule.BaseModule
	svr *service.Server
}

type MsgFormat struct {
	Module string
	Func   string
	Params interface{}
}

type resultInfo struct {
	Error     string  `json:"error,omitempty"` //错误结果 如果为nil表示请求正确
	Result    interface{} `json:"result"`      //结果
	ErrorCode int32    `json:"errorCode"`      //错误代码
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

	toResult := func(Topic []byte, Result interface{}, error *base.ErrorCode, packetId uint16) (err error) {
		r := &resultInfo{
			Error:     error.Desc,
			Result:    Result,
			ErrorCode: error.ErrorCode,
		}
		b, err := json.Marshal(r)
		if err == nil {
			m.WriteMsg(Topic, b, packetId)
		} else {
			r = &resultInfo{
				Error:  err.Error(),
				Result: nil,
			}
			log.Error(err.Error())

			br, _ := json.Marshal(r)
			m.WriteMsg(Topic, br, packetId)
		}
		return
	}
	log.Info(msg.String())
	topic := msg.Topic()
	payload := msg.Payload()
	var ArgsType []string = make([]string, 2)
	var args [][]byte = make([][]byte, 2)
	var msgContent MsgFormat
	packetId := msg.PacketId()
	err := json.Unmarshal(payload, &msgContent)
	if err == nil {
		arg1, err := json.Marshal(msgContent.Params)
		if err == nil {
			args[1] = arg1
		} else {
			log.Error("msg param format error %s", err.Error())
			toResult(topic, nil, base.ErrParamParseFail, packetId)
		}
		serverSession, err := m.App.GetRouteServers(msgContent.Module, "")
		if err != nil {
			log.Error("Service(type:%s) not found", msgContent.Module)
			toResult(topic, nil, base.ErrNotFound, packetId)
			return false
		}
		ArgsType[0] = RPC_PARAM_SESSION_TYPE

		b := []byte(sess.Id)
		args[0] = b
		result, e := serverSession.CallArgs(msgContent.Func, ArgsType, args)
		toResult(topic, result, e, packetId)
	} else {
		log.Error("msg format error %s", err.Error())
		toResult(topic, nil, base.ErrParamParseFail, packetId)
		return false
	}
	return true
}

func (m *Gate) WriteMsg(topic []byte, body []byte, packetId uint16) error {
	publish := message.NewPublishMessage()
	publish.SetTopic(topic)
	publish.SetPayload(body)
	publish.SetPacketId(packetId)
	return m.svr.Publish(publish, nil)
}

func (m *Gate) DisConnect(sess *sessions.Session) {
	log.Info("%s disconnect ", sess.Id)
}
