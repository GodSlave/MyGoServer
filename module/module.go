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
package module

import (
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/go-xorm/xorm"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/garyburd/redigo/redis"
)

type ServerSession interface {
	GetId() string
	GetType() string
	GetByteType() byte
	GetRpc() mqrpc.RPCClient
	CallArgs(_func string, sessionId string, args []byte) ([]byte, *base.ErrorCode)
	CallByteArgs(_func byte, sessionId string, args []byte) ([]byte, *base.ErrorCode)
	CallNRArgs(_func string, sessionId string, args []byte) (err error)
}
type App interface {
	Run(mods ...Module) error
	/**
	当同一个类型的Module存在多个服务时,需要根据情况选择最终路由到哪一个服务去
	fn: func(moduleType string,serverId string,[]*ServerSession)(*ServerSession)
	*/
	Route(moduleType string, fn func(app App, Type string, hash string) ServerSession) error
	Configure(settings conf.Config) error
	OnInit(settings conf.Config) error
	OnDestroy() error
	RegisterLocalClient(serverId string, server mqrpc.RPCServer) error
	GetServersById(id string) (ServerSession, error)
	/**
	filter		 调用者服务类型    moduleType|moduleType@moduleID
	Type	   	想要调用的服务类型
	*/
	GetRouteServers(filter string, hash string) (ServerSession, error)   //获取经过筛选过的服务
	GetByteRouteServers(filter byte, hash string) (ServerSession, error) //获取经过筛选过的服务
	GetServersByType(Type string) []ServerSession
	GetServersByByteType(Type byte) []ServerSession
	GetSettings() conf.Config //获取配置信息
	RpcInvokeNRArgs(module RPCModule, moduleType string, _func string, sessionId string, args []byte) (err error)
	RpcInvokeArgs(module RPCModule, moduleType string, _func string, sessionId string, args []byte) (result interface{}, err *base.ErrorCode)
	/**
	添加一个 自定义参数序列化接口
	gate,system 关键词一被占用请使用其他名称
	 */
	AddRPCSerialize(name string, Interface RPCSerialize) error

	GetRPCSerialize() (map[string]RPCSerialize)

	GetSqlEngine() *xorm.Engine

	GetRedis() *redis.Pool

	GetUserManager() UserManager
}
type ConnectEventCallBack func(sessionID string)

type Gate interface {
	SetOnConnectCallBack(callback ConnectEventCallBack)
	SetOnDisConnectCallBack(callback ConnectEventCallBack)
}

type UserEventCallBack func(user *base.BaseUser)

type UserManager interface {
	OnUserLogOut(user *base.BaseUser) //用于处理用户逻辑
	OnUserLogin(user *base.BaseUser)
	OnUserRegister(user *base.BaseUser)
	OnUserConnect(sessionId string)    //用户连接进来以后如果clientID是分配的token 表示是信任用户，为用户登陆
	OnUserDisconnect(sessionId string) //当用户断开连接后,清除用户的缓存
	VerifyUser(sessionId string) (user *base.BaseUser)  //根据sessionID 获取用户对象
	VerifyUserID(sessionId string) (userID string)  //根据sessionID获取用户
	SetLoginCallBack(callback UserEventCallBack)
	SetRegisterCallBack(callback UserEventCallBack)
	SetLogoutCallBack(callback UserEventCallBack)
}

type Module interface {
	Version() string                             //模块版本
	GetType() string                             //模块类型
	OnConfChanged(settings *conf.ModuleSettings) //为以后动态服务发现做准备
	OnInit(app App, settings *conf.ModuleSettings)
	OnDestroy()
	GetApp() (App)
	Run(closeSig chan bool)
}
type RPCModule interface {
	Module
	GetServerId() string //模块类型
	//RpcInvoke(moduleType string, _func string, params ...interface{}) (interface{}, *base.ErrorCode)
	//RpcInvokeNR(moduleType string, _func string, params ...interface{}) error
	RpcInvokeArgs(moduleType string, _func string, sessionID string, args []byte) (interface{}, *base.ErrorCode)
	RpcInvokeNRArgs(moduleType string, _func string, sessionID string, args []byte) error
	GetModuleSettings() (settings *conf.ModuleSettings)
	/**
	filter		 调用者服务类型    moduleType|moduleType@moduleID
	Type	   	想要调用的服务类型
	*/
	GetRouteServers(filter string, hash string) (ServerSession, error)
	GetStatistical() (statistical string, err error)
	GetExecuting() int64
}

/**
rpc 自定义参数序列化接口
 */
type RPCSerialize interface {
	/**
	序列化 结构体-->[]byte
	param 需要序列化的参数值
	@return ptype 当能够序列化这个值,并且正确解析为[]byte时 返回改值正确的类型,否则返回 ""即可
	@return p 解析成功得到的数据, 如果无法解析该类型,或者解析失败 返回nil即可
	@return err 无法解析该类型,或者解析失败 返回错误信息
	 */
	Serialize(param interface{}) (ptype string, p []byte, err error)
	/**
	反序列化 []byte-->结构体
	ptype 参数类型 与Serialize函数中ptype 对应
	b   参数的字节流
	@return param 解析成功得到的数据结构
	@return err 无法解析该类型,或者解析失败 返回错误信息
	 */
	Deserialize(ptype string, b []byte) (param interface{}, err error)
	/**
	返回这个接口能够处理的所有类型
	 */
	GetTypes() ([]string)
}
