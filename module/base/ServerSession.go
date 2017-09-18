// Copyright 2014 loolgame Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package basemodule

import (
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/base"
)

func NewServerSession(Id string, Stype string, Btype byte, Rpc mqrpc.RPCClient) (module.ServerSession) {
	session := &serverSession{
		Id:    Id,
		Stype: Stype,
		Rpc:   Rpc,
		Btype: Btype,
	}
	return session
}

type serverSession struct {
	Id    string
	Stype string
	Btype byte
	Rpc   mqrpc.RPCClient
}

func (c *serverSession) GetId() string {
	return c.Id
}
func (c *serverSession) GetType() string {
	return c.Stype
}
func (c *serverSession) GetRpc() mqrpc.RPCClient {
	return c.Rpc
}

func (c *serverSession) GetByteType() byte {
	return c.Btype
}

/**
消息请求 需要回复
*/
func (c *serverSession) CallArgs(_func string, sessionId string, args []byte) ([]byte, *base.ErrorCode) {
	return c.Rpc.CallArgs(_func, sessionId, args)
}

func (c *serverSession) CallByteArgs(_func byte, sessionId string, args []byte) ([]byte, *base.ErrorCode) {
	return c.Rpc.CallByteArgs(_func, sessionId, args)
}

/**
消息请求 bu需要回复
*/
func (c *serverSession) CallNRArgs(_func string, sessionId string, args []byte) (err error) {
	return c.Rpc.CallNRArgs(_func, sessionId, args)
}
