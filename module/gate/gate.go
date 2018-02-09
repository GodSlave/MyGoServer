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
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/utils"
)

var RPC_PARAM_SESSION_TYPE = "SESSION"

type Gate struct {
	basemodule.BaseModule
	svr             *service.Server
	redisPool       *redis.Pool
	connCallBack    module.ConnectEventCallBack
	disConnCallback module.ConnectEventCallBack
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
	m.redisPool = app.GetRedis()
	m.GetServer().RegisterGO("PushMessage", 1, m.PushMessage)
	m.GetServer().RegisterGO("KickOut", 2, m.KickOut)
}

func (m *Gate) Run(closeSig chan bool) {

	m.svr = &service.Server{
		KeepAlive:        15,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
		BusinessAgent:    m,
		Services:         utils.NewBeeMap(),
	}
	// Listen and serve connections at localhost:1883
	addr := m.GetModuleSettings().Settings["TCPAddr"].(string)
	m.svr.ListenAndServe(addr)
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

func (m *Gate) OnNewMessage(msg *message.PublishMessage, sess *sessions.Session) bool {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Gate  OnRecover error [%s]", r)
		}
	}()
	topic := msg.Topic()
	if topic[0] == 'i' {
		return m.progressJsonMessage(msg, sess)
	} else if topic[0] == 'f' {
		return m.progressProtoMessage(msg, sess)
	}

	return true
}

func (m *Gate) progressJsonMessage(msg *message.PublishMessage, sess *sessions.Session) bool {
	toResult := func(Topic []byte, Result []byte, error *base.ErrorCode, packetId uint16) (err error) {
		r := &AllResponse{
			Msg:    error.Desc,
			Result: Result,
			State:  error.ErrorCode & 0xff,
		}
		b, err := json.Marshal(r)
		if err == nil {
			err := m.WriteMsg(Topic, b, packetId, sess)
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			r = &AllResponse{
				Msg:    err.Error(),
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
		var err error
		payload, err = aesCipher.Decrypt(payload)
		log.Debug(string(payload))
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

		if string(arg1) == "null" {
			arg1 = nil
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
	toResult := func(Topic []byte, Result []byte, method [] byte, error *base.ErrorCode, packetId uint16, ) (err error) {
		r := &AllResponse{
			Msg:    error.Desc,
			Result: Result,
			State:  error.ErrorCode & 0xff,
		}

		if r.State != 0 {
			log.Info("meet a error %v", r.State)
		}

		b, err := proto.Marshal(r)
		if err == nil {
			realResult := make([]byte, len(b)+2)
			copy(realResult[0:2], method[0:2])
			copy(realResult[2:], b)
			err := m.WriteMsg(Topic, realResult, packetId, sess)
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			r = &AllResponse{
				Msg:    err.Error(),
				Result: nil,
			}
			log.Error(err.Error())
			br, _ := proto.Marshal(r)
			realResult := make([]byte, len(br)+2)
			copy(realResult[0:2], method[0:2])
			copy(realResult[2:], br)
			err := m.WriteMsg(Topic, realResult, packetId, sess)
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
		log.Debug("%v", payload)
		if err != nil {
			log.Error(err.Error())
			return false
		}
	}

	packetId := msg.PacketId()
	serverSession, err := m.App.GetByteRouteServers(payload[0], "")
	if err != nil {
		log.Error("Service(type:%x) not found", payload[0])
		toResult(topic, nil, payload[0:2], base.ErrNotFound, packetId)
		return false
	}

	result, e := serverSession.CallByteArgs(payload[1], sess.Id, payload[2:])
	toResult(topic, result, payload[0:2], e, packetId)
	return true
}

func (m *Gate) WriteMsg(topic []byte, body []byte, packetId uint16, sess *sessions.Session) error {
	publish := message.NewPublishMessage()
	publish.SetTopic(topic)
	if m.GetApp().GetSettings().Secret {
		aesCipher, _ := aes.NewAesEncrypt(sess.AesKey)
		var err error
		body, err = aesCipher.EncryptBytes(body)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	publish.SetPayload(body)
	publish.SetQoS(1)
	publish.SetPacketId(packetId)
	return m.svr.PublishToClient(publish, sess.Id, nil)
}

func (m *Gate) OnDisConnect(sess *sessions.Session) {
	log.Debug("%s disconnect ", sess.Id)
	if m.disConnCallback != nil {
		m.disConnCallback(sess.Id)
	}
	//m.App.OnUserLogOut(sess.Id)
	if m.App.GetUserManager() != nil {
		m.App.GetUserManager().OnUserDisconnect(sess.Id)
	}
}

func (m *Gate) OnConnect(sess *sessions.Session) {
	if m.connCallBack != nil {
		m.connCallBack(sess.Id)
	}

	if m.App.GetUserManager() != nil {
		m.App.GetUserManager().OnUserConnect(sess.Id)
	}
}

func (m *Gate) SetOnConnectCallBack(callback module.ConnectEventCallBack) {
	m.connCallBack = callback
}
func (m *Gate) SetOnDisConnectCallBack(callback module.ConnectEventCallBack) {
	m.disConnCallback = callback
}

func (m *Gate) PushMessage(userId string, item *base.PushItem) {
	service := m.getService(userId)
	if service != nil {
		topic := []byte{'p'}
		content, err := json.Marshal(item)
		if err == nil {
			m.WriteMsg(topic, content, 0, service.GetSession())
		} else {
			log.Error(err.Error())
		}
		log.Debug(string(content))

		topic = []byte{'q'}
		topic2 := []byte{'f'}
		protoContent, _ := proto.Marshal(interface{}(item.Content).(proto.Message))
		r := &AllResponse{
			Msg:    "",
			Result: protoContent,
			State:  0,
		}

		realProtoMsg, _ := proto.Marshal(interface{}(r).(proto.Message))

		if err == nil {
			realContent := make([]byte, len(content)+2)
			realContent[0] = item.Module
			realContent[1] = item.PushType
			copy(realContent[2:], content)
			m.WriteMsg(topic, realProtoMsg, 0, service.GetSession())
			m.WriteMsg(topic2, realProtoMsg, 0, service.GetSession())
		} else {
			log.Error(err.Error())
		}
	} else {
	}

}

func (m *Gate) getService(userId string) *service.Service {
	if obj := m.svr.Services.Get(userId); obj != nil {
		return obj.(*service.Service)
	} else {
		conn, err := m.redisPool.Dial()
		if err == nil {
			var token, session string
			session, _ = redis.String(conn.Do("GET", base.ID_SESSION_PREFIX+userId))
			token, _ = redis.String(conn.Do("GET", base.ID_TOKEN_PERFIX+userId))

			if session != "" {
				if obj := m.svr.Services.Get(session); obj != nil {
					return obj.(*service.Service)
				}
			}

			if token != "" {
				if obj := m.svr.Services.Get(token); obj != nil {
					return obj.(*service.Service)
				}
			}
		}
	}
	log.Info("userModule Client not Connect :%s", userId)
	return nil
}

func (m *Gate) KickOut(userId string) {
	if obj := m.svr.Services.Get(userId); obj != nil {
		service := obj.(*service.Service)
		if service != nil {
			service.KickOut()
		}
	}
}

func (m *Gate) GetExecuting() int64 {
	if m.svr != nil {
		return m.svr.GetConnectUser()
	}
	return 0
}
