package master

import (
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/module"
)

var (
	Register        = "Register"
	OnRegister      = "OnRegister"
	UnRegister      = "UnRegister"
	OnUnRegister    = "OnUnRegister"
	GetAppStatus    = "GetAppStatus"
	ReportStatus    = "ReportStatus"
	SyncVersionCode = "SyncVersionCode"
	getAppList      = "getAppList"
	MasterStr       = "Master"
	UpdateStatus    = "UpdateStatus"
	Bye             = "Bye"
)

type Master interface {
	Register(info ApplicationInfo)        //服务注册到master
	UnRegister(appName string)            //解除注册
	GetAppStatus(name string) []AppStatus // 获取可用服务状态
	GetAppList(name string) []ApplicationInfo
	ReportStatus(status AppStatus) // 报告状态
}

type MasterClient interface {
	SetShutDownCallBack(callback ShutDownCallBack)                           //注册远程关闭回调
	SetUpdateApplicationListCallBack(callback UpdateApplicationListCallBack) //注册更新服务器列表回调
	GetModule(module string) *module.ServerSession                           //获取可用模块RPC信息
	GetModuleByByte(appByteName byte) *module.ServerSession                  //根据byte编号获取可用模块RPC信息
	Shutdown()                                                               //关闭服务
	ToShutdown()                                                             //主服务器关闭本服务
}

type ApplicationInfo struct {
	Name    string
	RpcInfo string
	Modules map[string][]*conf.ModuleSettings
}

type ShutDownCallBack func()              //关闭某个服务
type UpdateApplicationListCallBack func() // 更新服务器列表

type AppStatus struct {
	AppName           string
	Load              int32 //用0到100 表示负载高低
	ModuleStatus      []ModuleStatus
	ProcessIngNumbers int32 //正在处理中的请求数量
}

type ModuleStatus struct {
	ModuleName  string
	Load        int32  //当前模块的使用频率
	MethodLoads string //方法的调用信息
}

type OtherModuleInfo struct {
	serverSession *module.ServerSession
	load          int
	appName       string
	key           string
}
