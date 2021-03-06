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
	"fmt"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/golang/protobuf/proto"
	"github.com/GodSlave/MyGoServer/utils"
	"sync"
	"time"
	"github.com/GodSlave/MyGoServer/rpc/pb"
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/GodSlave/MyGoServer/rpc/util"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/base"
	"container/list"
)

type RedisClient struct {
	//callInfos map[string]*ClinetCallInfo
	callInfos         *utils.BeeMap
	cmutex            sync.Mutex //操作
	callLog           *list.List
	psc               redis.PubSubConn
	info              *conf.Redis
	queueName         string
	callbackQueueName string
	done              chan error
	ticker            *time.Ticker
	timeoutResponse   *rpcpb.ResultInfo
}

type CallLogInfo struct {
	time int64  //放入时间
	cid  string //请求id
}

func createQueueName() string {
	//return "callbackQueueName"
	return fmt.Sprintf("callbackQueueName:%d", time.Now().Nanosecond())
}
func NewRedisClient(info *conf.Redis) (client *RedisClient, err error) {
	var url = info.Uri
	client = new(RedisClient)
	client.callInfos = utils.NewBeeMap()
	client.info = info
	client.callbackQueueName = createQueueName()
	client.queueName = info.Queue
	client.done = make(chan error)
	psc := redis.PubSubConn{Conn: utils.GetRedisFactory().GetPool(url).Get()}
	psc.Subscribe(client.callbackQueueName)

	pool := utils.GetRedisFactory().GetPool(info.Uri).Get()
	defer pool.Close()
	_, errs := pool.Do("EXPIRE", client.callbackQueueName, 60)
	if errs != nil {
		log.Warning(errs.Error())
	}
	client.psc = psc
	go client.on_response_handle(client.done)
	client.ticker = time.NewTicker(time.Second * 1)
	client.callLog = list.New()
	client.timeoutResponse = &rpcpb.ResultInfo{
		Result:     nil,
		Error:      base.ErrRequestTimeout.Desc,
		ErrorCode:  base.ErrRequestTimeout.ErrorCode,
		ResultType: argsutil.NULL,
	}
	client.handTimeOut()
	return client, nil
}

func (c *RedisClient) Done() (err error) {
	pool := utils.GetRedisFactory().GetPool(c.info.Uri).Get()
	defer pool.Close()
	//关闭amqp链接通道
	c.psc.Unsubscribe(c.callbackQueueName)
	//删除临时通道
	pool.Do("DEL", c.callbackQueueName)
	//err = c.psc.Close()
	//清理 callInfos 列表
	for key, clientCallInfo := range c.callInfos.Items() {
		if clientCallInfo != nil {
			//关闭管道
			close(clientCallInfo.(ClinetCallInfo).call)
			//从Map中删除
			c.callInfos.Delete(key)
		}
	}
	c.callInfos = nil
	if c.ticker != nil {
		c.ticker.Stop()
		log.Error("stopped ticker")
	}
	return
}

/**
消息请求
*/
func (c *RedisClient) Call(callInfo mqrpc.CallInfo, callback chan rpcpb.ResultInfo) error {
	pool := utils.GetRedisFactory().GetPool(c.info.Uri).Get()
	defer pool.Close()
	var err error
	if c.callInfos == nil {
		return fmt.Errorf("RedisClient is closed")
	}
	callInfo.RpcInfo.ReplyTo = c.callbackQueueName
	var correlation_id = callInfo.RpcInfo.Cid
	clientCallInfo := &ClinetCallInfo{
		correlation_id: correlation_id,
		call:           callback,
		timeout:        callInfo.RpcInfo.Expired,
	}
	c.callInfos.Set(correlation_id, *clientCallInfo)
	c.callLog.PushBack(&CallLogInfo{
		time: callInfo.RpcInfo.Expired,
		cid:  correlation_id,
	})
	body, err := c.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}

	_, err = pool.Do("PUBLISH", c.queueName, body)
	if err != nil {
		log.Warning("Publish: %s", err)
		return err
	}
	return nil
}

/**
消息请求 不需要回复
*/
func (c *RedisClient) CallNR(callInfo mqrpc.CallInfo) error {
	pool := utils.GetRedisFactory().GetPool(c.info.Uri).Get()
	defer pool.Close()
	var err error

	body, err := c.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	_, err = pool.Do("PUBLISH", c.queueName, body)
	if err != nil {
		log.Warning("Publish: %s", err)
		return err
	}
	return nil
}

/**
接收应答信息
*/
func (c *RedisClient) on_response_handle(done chan error) {
	for {
		switch v := c.psc.Receive().(type) {
		case redis.Message:
			resultInfo, err := c.UnmarshalResult(v.Data)
			if err != nil {
				log.Error("Unmarshal faild", err)
			} else {
				correlation_id := resultInfo.Cid
				clinetCallInfo := c.callInfos.Get(correlation_id)
				if clinetCallInfo != nil {
					clinetCallInfo.(ClinetCallInfo).call <- *resultInfo
				}
				//删除
				//log.Debug("remove %s", correlation_id)
				c.callInfos.Delete(correlation_id)
			}
		case redis.PMessage:
			resultInfo, err := c.UnmarshalResult(v.Data)
			if err != nil {
				log.Error("Unmarshal faild", err)
			} else {
				correlation_id := resultInfo.Cid
				clinetCallInfo := c.callInfos.Get(correlation_id)
				if clinetCallInfo != nil {
					clinetCallInfo.(ClinetCallInfo).call <- *resultInfo
				}
				//删除
				log.Info("remove %s", correlation_id)
				c.callInfos.Delete(correlation_id)
			}
		case redis.Subscription:
			log.Info("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Error("on_response_handle", v.Error())
			return
		default:

		}

	}
	log.Debug("finish on_response_handle")
}

func (c *RedisClient) handTimeOut() {
	go func() {
		for {
			select {
			case <-c.ticker.C:
			Skip:
				for {
					first := c.callLog.Front()
					if first != nil {
						callInfo := first.Value.(*CallLogInfo)
						if time.Now().Unix() > callInfo.time {
							c.callLog.Remove(first)
							//发送一个超时的消息
							value := c.callInfos.Get(callInfo.cid)
							if value != nil {
								var clientCallInfo = value.(ClinetCallInfo)
								if clientCallInfo.call != nil {
									clientCallInfo.call <- *c.timeoutResponse
								}
								c.callInfos.Delete(callInfo.cid)
							}
						} else {
							break Skip
						}
					} else {
						break Skip
					}
				}
			}
		}
	}()
}

func (c *RedisClient) UnmarshalResult(data []byte) (*rpcpb.ResultInfo, error) {
	//fmt.Println(msg)
	//保存解码后的数据，Value可以为任意数据类型
	var resultInfo rpcpb.ResultInfo
	err := proto.Unmarshal(data, &resultInfo)
	if err != nil {
		return nil, err
	} else {
		return &resultInfo, err
	}
}

func (c *RedisClient) Unmarshal(data []byte) (*rpcpb.RPCInfo, error) {
	//fmt.Println(msg)
	//保存解码后的数据，Value可以为任意数据类型
	var rpcInfo rpcpb.RPCInfo
	err := proto.Unmarshal(data, &rpcInfo)
	if err != nil {
		return nil, err
	} else {
		return &rpcInfo, err
	}
	panic("bug")
}

// goroutine safe
func (c *RedisClient) Marshal(rpcInfo *rpcpb.RPCInfo) ([]byte, error) {
	//map2:= structs.Map(callInfo)
	b, err := proto.Marshal(rpcInfo)
	return b, err
}
