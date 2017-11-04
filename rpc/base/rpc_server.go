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
	"github.com/GodSlave/MyGoServer/rpc/pb"
	"reflect"
	"sync"
	"time"
	"runtime"
	"github.com/GodSlave/MyGoServer/rpc/util"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/opentracing/opentracing-go"
	"encoding/json"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/gogo/protobuf/proto"
	serverbase "github.com/GodSlave/MyGoServer/base"
)

type RPCServer struct {
	module         module.Module
	app            module.App
	functions      map[string]mqrpc.FunctionInfo
	byteFunctions  map[int32]mqrpc.FunctionInfo
	remote_server  *AMQPServer
	local_server   *LocalServer
	redis_server   *RedisServer
	mq_chan        chan mqrpc.CallInfo //接收到请求信息的队列
	callback_chan  chan mqrpc.CallInfo //信息处理完成的队列
	wg             sync.WaitGroup      //任务阻塞
	call_chan_done chan error
	listener       mqrpc.RPCListener
	executing      int64 //正在执行的goroutine数量
}

func NewRPCServer(app module.App, module module.Module) (mqrpc.RPCServer, error) {
	rpc_server := new(RPCServer)
	rpc_server.app = app
	rpc_server.module = module
	rpc_server.call_chan_done = make(chan error)
	rpc_server.functions = make(map[string]mqrpc.FunctionInfo)
	rpc_server.byteFunctions = make(map[int32]mqrpc.FunctionInfo)
	rpc_server.mq_chan = make(chan mqrpc.CallInfo, 50)
	rpc_server.callback_chan = make(chan mqrpc.CallInfo, 50)

	//先创建一个本地的RPC服务
	local_server, err := NewLocalServer(rpc_server.mq_chan)
	if err != nil {
		log.Error("LocalServer Dial: %s", err)
	}
	rpc_server.local_server = local_server

	go rpc_server.on_call_handle(rpc_server.mq_chan, rpc_server.callback_chan, rpc_server.call_chan_done)

	go rpc_server.on_callback_handle(rpc_server.callback_chan) //结果发送队列
	return rpc_server, nil
}

/**
创建一个支持远程RPC的服务
*/
func (s *RPCServer) NewRabbitmqRPCServer(info *conf.Rabbitmq) (err error) {
	remote_server, err := NewAMQPServer(info, s.mq_chan)
	if err != nil {
		log.Error("AMQPServer Dial: %s", err)
	}
	s.remote_server = remote_server
	return
}

/**
创建一个支持远程Redis RPC的服务
*/
func (s *RPCServer) NewRedisRPCServer(info *conf.Redis) (err error) {
	redis_server, err := NewRedisServer(info, s.mq_chan)
	if err != nil {
		log.Error("RedisServer Dial: %s", err)
	}
	s.redis_server = redis_server
	return
}
func (s *RPCServer) SetListener(listener mqrpc.RPCListener) {
	s.listener = listener
}
func (s *RPCServer) GetLocalServer() mqrpc.LocalServer {
	return s.local_server
}

/**
获取当前正在执行的goroutine 数量
*/
func (s *RPCServer) GetExecuting() int64 {
	return s.executing
}

// you must call the function before calling Open and Go
func (s *RPCServer) Register(id string, byteId int32, f interface{}) {

	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	fuction1 := *&mqrpc.FunctionInfo{
		Function:  f,
		Goroutine: false,
	}

	s.functions[id] = fuction1
	s.byteFunctions[byteId] = fuction1
}

// you must call the function before calling Open and Go
func (s *RPCServer) RegisterGO(id string, byteId int32, f interface{}) {

	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	fuction1 := *&mqrpc.FunctionInfo{
		Function:  f,
		Goroutine: true,
	}
	s.functions[id] = fuction1
	s.byteFunctions[byteId] = fuction1
}

func (s *RPCServer) Done() (err error) {
	//设置队列停止接收请求
	if s.remote_server != nil {
		err = s.remote_server.StopConsume()
	}
	if s.local_server != nil {
		err = s.local_server.StopConsume()
	}
	//等待正在执行的请求完成
	close(s.mq_chan)   //关闭mq_chan通道
	<-s.call_chan_done //mq_chan通道的信息都已处理完
	s.wg.Wait()
	close(s.callback_chan) //关闭结果发送队列
	//关闭队列链接
	if s.remote_server != nil {
		err = s.remote_server.Shutdown()
	}
	if s.local_server != nil {
		err = s.local_server.Shutdown()
	}
	return
}

/**
处理结果信息
*/
func (s *RPCServer) on_callback_handle(callbacks <-chan mqrpc.CallInfo) {
	for {
		select {
		case callInfo, ok := <-callbacks:
			if !ok {
				callbacks = nil
			} else {
				if callInfo.RpcInfo.Reply {
					//需要回复的才回复
					callInfo.Agent.(mqrpc.MQServer).Callback(callInfo)
				} else {
					//对于不需要回复的消息,可以判断一下是否出现错误，打印一些警告
					if callInfo.Result.Error != "" {
						log.Warning("rpc callback erro :\n%s", callInfo.Result.Error)
					}
				}
			}
		}
		if callbacks == nil {
			break
		}
	}
}

/**
接收请求信息
*/
func (s *RPCServer) on_call_handle(calls <-chan mqrpc.CallInfo, callbacks chan<- mqrpc.CallInfo, done chan error) {
	for {
		select {
		case callInfo, ok := <-calls:
			if !ok {
				calls = nil
			} else {
				if callInfo.RpcInfo.Expired < (time.Now().UnixNano() / 1000000) {
					//请求超时了,无需再处理
					if s.listener != nil {
						s.listener.OnTimeOut(callInfo.RpcInfo.Fn, callInfo.RpcInfo.ByteFn, callInfo.RpcInfo.Expired)
					} else {
						log.Warning("timeout: This is Call", s.module.GetType(), callInfo.RpcInfo.Fn, callInfo.RpcInfo.Expired, time.Now().UnixNano()/1000000)
					}
				} else {
					s.runFunc(callInfo, callbacks)
				}
			}
		}
		if calls == nil {
			done <- nil
			break
		}
	}
}

//---------------------------------if _func is not a function or para num and type not match,it will cause panic
func (s *RPCServer) runFunc(callInfo mqrpc.CallInfo, callbacks chan<- mqrpc.CallInfo) {
	_errorCallback := func(Cid string, errorCode *base.ErrorCode) {
		resultInfo := rpcpb.NewResultInfo(Cid, errorCode.Desc, errorCode.ErrorCode, argsutil.NULL, nil)
		callInfo.Result = *resultInfo
		callbacks <- callInfo
		if s.listener != nil {
			s.listener.OnError(callInfo.RpcInfo.Fn, &callInfo, fmt.Errorf(errorCode.Error()))
		}
	}
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			error := base.NewError(500, rn)
			_errorCallback(callInfo.RpcInfo.Cid, error)
		}
	}()

	var functionInfo mqrpc.FunctionInfo
	var ok bool
	if callInfo.RpcInfo.Fn != "" {
		functionInfo, ok = s.functions[callInfo.RpcInfo.Fn]
	} else {
		functionInfo, ok = s.byteFunctions[callInfo.RpcInfo.ByteFn]
	}

	if !ok {
		_errorCallback(callInfo.RpcInfo.Cid, base.NewError(404, fmt.Sprintf("Remote function(%s) not found", callInfo.RpcInfo.Fn)))
		return
	}
	_func := functionInfo.Function

	params := callInfo.RpcInfo.Args
	//ArgsType := callInfo.RpcInfo.ArgsType
	f := reflect.ValueOf(_func)
	funcType := reflect.TypeOf(_func)
	//if len(params) != f.Type().NumIn() {
	//	//因为在调研的 _func的时候还会额外传递一个回调函数 cb
	//	_errorCallback(callInfo.RpcInfo.Cid, fmt.Sprintf("The number of params %s is not adapted.%s", params, f.String()))
	//	return
	//}

	s.wg.Add(1)
	s.executing++
	_runFunc := func() {
		var span opentracing.Span = nil

		defer func() {
			if r := recover(); r != nil {
				var rn *serverbase.ErrorCode
				switch r.(type) {
				case *serverbase.ErrorCode:
					rn = r.(*serverbase.ErrorCode)
				case string:
					str := r.(string)
					rn = base.NewError(500, str)
				case error:
					str := r.(error).Error()
					rn = base.NewError(500, str)
				}
				buf := make([]byte, 1024)
				l := runtime.Stack(buf, false)
				errstr := string(buf[:l])
				log.Error("%s rpc func(%s) error %s\n ----Stack----\n%s", s.module.GetType(), callInfo.RpcInfo.Fn, rn.Desc, errstr)
				_errorCallback(callInfo.RpcInfo.Cid, rn)
			}

			if span != nil {
				span.Finish()
			}

			s.wg.Add(-1)
			s.executing--
		}()
		exec_time := time.Now().UnixNano()
		//t:=RandInt64(2,3)
		//time.Sleep(time.Second*time.Duration(t))
		// f 为函数地址
		var session string = ""
		var in []reflect.Value
		in = make([]reflect.Value, funcType.NumIn())
		userIndex := 0
		paramsIndex := 0
		if len(in) > 1 {
			userIndex = 0
			paramsIndex = 1
		} else if len(in) == 1 {
			if params == nil || len(params) == 0 {
				userIndex = 0
				paramsIndex = -1
			} else {
				userIndex = -1
				paramsIndex = 0
			}
		}

		if userIndex >= 0 {
			param1type := funcType.In(userIndex)
			//TODO map session to real user id
			sessionID := callInfo.RpcInfo.SessionId
			if param1type.Kind() == reflect.String {
				userID := s.app.GetUserManager().VerifyUserID(sessionID)
				if userID != "" {
					in[userIndex] = reflect.ValueOf(userID)
				} else {
					in[userIndex] = reflect.ValueOf(sessionID)
				}
			} else {
				user := s.app.GetUserManager().VerifyUser(sessionID)
				if user != nil {
					in[userIndex] = reflect.ValueOf(user)
				} else {
					_errorCallback(callInfo.RpcInfo.Cid, base.ErrNeedLogin)
					return
				}
			}
		}
		if paramsIndex >= 0 {
			paramType := funcType.In(paramsIndex)
			v := reflect.New(paramType.Elem()).Interface()
			var err error
			if callInfo.RpcInfo.Fn != "" {
				err = json.Unmarshal(params, &v)
			} else {
				err = proto.Unmarshal(params, v.(proto.Message))
			}
			if err != nil {
				log.Error(err.Error())
				panic(err)
			}
			in[paramsIndex] = reflect.ValueOf(v)
		}

		if s.listener != nil {
			errs := s.listener.BeforeHandle(callInfo.RpcInfo.Fn, session, &callInfo)
			if errs != nil {
				_errorCallback(callInfo.RpcInfo.Cid, base.ErrInternal)
				return
			}
		}
		out := f.Call(in)
		errorCode := base.ErrNil
		var result []byte

		for k, value := range out {
			if k == 0 {
				switch value.Kind() {
				case reflect.Invalid:
					result = nil
					log.Info("invalid result ")
				case reflect.String:
					result = []byte(value.Interface().(string))
				default:
					if callInfo.RpcInfo.Fn != "" {
						result, _ = json.Marshal(value.Interface())
					} else {
						result, _ = proto.Marshal(value.Interface().(proto.Message))
					}
				}
			}

			if k == 1 {
				if value.IsNil() {
					errorCode = base.ErrNil
				} else {
					switch value.Kind() {
					case reflect.Invalid:
						errorCode = base.ErrNil
					case reflect.String:
						errorInfo := value.String()
						if errorInfo != "" {
							errorCode = base.NewError(500, errorInfo)
						}
					default:
						errorCode1, err := value.Interface().(*base.ErrorCode)
						if err {
							errorCode = errorCode1
						}
					}
				}

			}
		}

		if len(out) < 2 {
			_errorCallback(callInfo.RpcInfo.Cid, base.ErrInternal)
			return
		}

		resultInfo := rpcpb.NewResultInfo(
			callInfo.RpcInfo.Cid,
			errorCode.Desc,
			errorCode.ErrorCode,
			argsutil.BYTES,
			result,
		)
		callInfo.Result = *resultInfo
		callbacks <- callInfo

		if s.listener != nil {
			s.listener.OnComplete(callInfo.RpcInfo.Fn, &callInfo, resultInfo, time.Now().UnixNano()-exec_time)
		}
	}
	if functionInfo.Goroutine {
		go _runFunc()
	} else {
		_runFunc()
	}
}
