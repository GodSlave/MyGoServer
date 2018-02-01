package master

import (
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/module"
)

var (
	Register        = "Register"        //子服务器注册到主服务器
	OnRegister      = "OnRegister"      //主服务器通知各个子服务器
	UnRegister      = "UnRegister"      //子服务器通知到主服务器
	OnUnRegister    = "OnUnRegister"    //主服务器通知各个子服务器
	GetAppStatus    = "GetAppStatus"    //获得各个子服务器状态
	ReportStatus    = "ReportStatus"    //子服务器报告自身状态
	SyncVersionCode = "SyncVersionCode" //同步版本号
	GetAppList      = "GetAppList"      //获得子服务器列表
	MasterStr       = "Master"          //主服务器标志
	UpdateStatus    = "UpdateStatus"    //主服务器同步子服务器状态
	Bye             = "Bye"
	UpdateInfo      = "UpdateInfo" //主服务器要求子服务器重新注册·
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
	load          float32
	appName       string
	key           string
}
