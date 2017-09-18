// Copyright 2014 mqant Author. All Rights Reserved.
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
package defaultrpc

import (
	"time"
	"fmt"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"github.com/GodSlave/MyGoServer/rpc/pb"
	"github.com/golang/protobuf/proto"
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/base"
)

type RPCClient struct {
	app           module.App
	serverId      string
	remote_client *AMQPClient
	local_client  *LocalClient
	redis_client  *RedisClient
}

func NewRPCClient(app module.App, serverId string) (mqrpc.RPCClient, error) {
	rpc_client := new(RPCClient)
	rpc_client.serverId = serverId
	rpc_client.app = app
	return rpc_client, nil
}

func (c *RPCClient) NewRabbitmqClient(info *conf.Rabbitmq) (err error) {
	//创建本地连接
	if info != nil && c.remote_client == nil {
		c.remote_client, err = NewAMQPClient(info)
		if err != nil {
			log.Error("Dial: %s", err)
		}
	}
	return
}

func (c *RPCClient) NewLocalClient(server mqrpc.RPCServer) (err error) {
	//创建本地连接
	if server != nil && server.GetLocalServer() != nil && c.local_client == nil {
		c.local_client, err = NewLocalClient(server.GetLocalServer())
		if err != nil {
			log.Error("Dial: %s", err)
		}
	}
	return
}

func (c *RPCClient) NewRedisClient(info *conf.Redis) (err error) {
	//创建本地连接
	if info != nil && c.redis_client == nil {
		c.redis_client, err = NewRedisClient(info)
		if err != nil {
			log.Error("Dial: %s", err)
		}
	}
	return
}

func (c *RPCClient) Done() (err error) {
	if c.remote_client != nil {
		err = c.remote_client.Done()
	}
	if c.redis_client != nil {
		err = c.redis_client.Done()
	}
	if c.local_client != nil {
		err = c.local_client.Done()
	}
	return
}

func (c *RPCClient) CallArgs(_func string, sessionId string, args []byte) ([]byte, *base.ErrorCode) {
	var correlation_id = uuid.Rand().Hex()
	rpcInfo := &rpcpb.RPCInfo{
		Fn:        *proto.String(_func),
		Reply:     *proto.Bool(true),
		Expired:   *proto.Int64((time.Now().UTC().Add(time.Second * time.Duration(c.app.GetSettings().Rpc.RpcExpired)).UnixNano()) / 1000000),
		Cid:       *proto.String(correlation_id),
		SessionId: *proto.String(sessionId),
		Args:      args,
	}
	return c.processCallInfo(rpcInfo)
}

func (c *RPCClient) CallByteArgs(_func byte, SessionId string, args []byte) ([]byte, *base.ErrorCode) {
	var correlation_id = uuid.Rand().Hex()
	rpcInfo := &rpcpb.RPCInfo{
		ByteFn:    *proto.Int32(int32(_func)),
		Reply:     *proto.Bool(true),
		Expired:   *proto.Int64((time.Now().UTC().Add(time.Second * time.Duration(c.app.GetSettings().Rpc.RpcExpired)).UnixNano()) / 1000000),
		Cid:       *proto.String(correlation_id),
		SessionId: *proto.String(SessionId),
		Args:      args,
	}
	return c.processCallInfo(rpcInfo)
}

func (c *RPCClient) processCallInfo(rpcInfo *rpcpb.RPCInfo) ([]byte, *base.ErrorCode) {
	callInfo := &mqrpc.CallInfo{
		RpcInfo: *rpcInfo,
	}
	callback := make(chan rpcpb.ResultInfo, 1)
	var err error
	//优先使用本地rpc
	if c.local_client != nil {
		err = c.local_client.Call(*callInfo, callback)
	} else if c.remote_client != nil {
		err = c.remote_client.Call(*callInfo, callback)
	} else if c.redis_client != nil {
		err = c.redis_client.Call(*callInfo, callback)
	} else {
		log.Error("rpc service (%s) connection failed", c.serverId)
		return nil, base.ErrServerIsDown
	}

	resultInfo, ok := <-callback
	if !ok {
		log.Error("client closed")
		return nil, base.ErrServerIsDown
	}

	if err != nil {
		log.Error(err.Error())
		return nil, base.ErrInternal
	}
	return resultInfo.Result, base.NewError(resultInfo.ErrorCode, resultInfo.Error)
}

func (c *RPCClient) CallNRArgs(_func string, SessionID string, args []byte) (err error) {
	var correlation_id = uuid.Rand().Hex()
	rpcInfo := &rpcpb.RPCInfo{
		Fn:        *proto.String(_func),
		Reply:     *proto.Bool(false),
		Expired:   *proto.Int64((time.Now().UTC().Add(time.Second * time.Duration(c.app.GetSettings().Rpc.RpcExpired)).UnixNano()) / 1000000),
		Cid:       *proto.String(correlation_id),
		Args:      args,
		SessionId: SessionID,
	}
	callInfo := &mqrpc.CallInfo{
		RpcInfo: *rpcInfo,
	}
	//优先使用本地rpc
	if c.local_client != nil {
		err = c.local_client.CallNR(*callInfo)
	} else if c.remote_client != nil {
		err = c.remote_client.CallNR(*callInfo)
	} else if c.redis_client != nil {
		err = c.redis_client.CallNR(*callInfo)
	} else {
		return fmt.Errorf("rpc service (%s) connection failed", c.serverId)
	}
	return nil
}
