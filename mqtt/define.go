package mymqtt

import (
	"github.com/surgemq/message"
	"github.com/GodSlave/MyGoServer/mqtt/sessions"
)

type MsgProcess interface {
	OnNewMessage(msg *message.PublishMessage, sess *sessions.Session) bool

	OnDisConnect(sess *sessions.Session)

	OnConnect(sess *sessions.Session)

}
