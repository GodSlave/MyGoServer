package mymqtt

import (
	"github.com/surgemq/message"
	"github.com/GodSlave/MyGoServer/mqtt/sessions"
)

type MsgProcess interface {
	Process(msg *message.PublishMessage, sess *sessions.Session) bool

	DisConnect(sess *sessions.Session)

}
