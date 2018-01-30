package master

import (
	"github.com/GodSlave/MyGoServer/rpc/base"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/GodSlave/MyGoServer/rpc/pb"
	"strconv"
	"time"
	"encoding/json"
	"github.com/GodSlave/MyGoServer/log"
	"sync"
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/module"
)

type DefaultMasterClient struct {
	MasterClient
	shutDownCallBack ShutDownCallBack
	updateCallBack   UpdateApplicationListCallBack
	rpcClient        *defaultrpc.RedisClient
	rpcServer        *defaultrpc.RedisServer
	selfRpcServer    *defaultrpc.RedisServer
	moduleInfo       map[string][]*OtherModuleInfo
	moduleInfoByte   map[byte][]*OtherModuleInfo
	Name             string
	callId           int
	callback_chan    chan rpcpb.ResultInfo
	versionCode      int32
	app              module.App
	lock             *sync.RWMutex
	moduleManager    *basemodule.ModuleManager
}

func NewMasterClient(config conf.Master, appName string, app module.App, moduleManage *basemodule.ModuleManager) DefaultMasterClient {
	m := DefaultMasterClient{}
	m.app = app
	m.moduleManager = moduleManage
	m.lock = new(sync.RWMutex)
	var err error
	m.rpcClient, err = defaultrpc.NewRedisClient(config.RedisPubSubConf)
	if err != nil {
		panic(err)
	}

	call_chan := make(chan mqrpc.CallInfo, 10)

	m.rpcServer, err = defaultrpc.NewRedisServer(config.RedisPubSubConf, call_chan)
	serverRedisConf := &conf.Redis{
		Uri:   config.RedisPubSubConf.Uri,
		Queue: appName,
	}

	self_call_chan := make(chan mqrpc.CallInfo, 10)
	m.selfRpcServer, err = defaultrpc.NewRedisServer(serverRedisConf, self_call_chan)
	m.Name = appName
	m.moduleInfo = map[string][]*OtherModuleInfo{}
	m.moduleInfoByte = map[byte][]*OtherModuleInfo{}
	go m.startListen(call_chan, self_call_chan)
	go m.RegisterToServer()
	go m.tick()
	return m
}

func (m DefaultMasterClient) SetShutDownCallBack(callback ShutDownCallBack) {
	m.shutDownCallBack = callback
}

func (m DefaultMasterClient) SetUpdateApplicationListCallBack(callback UpdateApplicationListCallBack) {
	m.updateCallBack = callback
}

func (m DefaultMasterClient) GetModule(module string) *module.ServerSession {
	modules := m.moduleInfo[module]
	if modules != nil && len(modules) > 0 {
		return modules[0].serverSession
	}

	return nil
}

func (m DefaultMasterClient) GetModuleByByte(appByteName byte) *module.ServerSession {
	modules := m.moduleInfoByte[appByteName]
	if modules != nil && len(modules) > 0 {
		return modules[0].serverSession
	}
	return nil
}

func (m DefaultMasterClient) Shutdown() {
	//TODO
	m.publicMessage(Bye, m.Name, Bye)

}

func (m DefaultMasterClient) ToShutdown() {
	m.shutDownCallBack()
	//TODO
}

func (m DefaultMasterClient) RegisterToServer() {
	time.Sleep(1 * time.Second)
	appInfo := &ApplicationInfo{
		Name:    m.Name,
		Modules: conf.Conf.Module,
	}
	m.publicMessage(Register, m.Name, appInfo)
	m.publicMessage(getAppList, m.Name, m.Name)
}

func (m DefaultMasterClient) buildDefaultCallInfo(functionName string, from string, args []byte) *mqrpc.CallInfo {
	m.callId += 1
	callInfo := mqrpc.CallInfo{
		RpcInfo: rpcpb.RPCInfo{
			Cid:       m.Name + strconv.Itoa(m.callId),
			Fn:        functionName,
			Args:      args,
			Reply:     false,
			SessionId: from,
			ByteFn:    m.versionCode,
		},
	}
	return &callInfo;
}

func (m DefaultMasterClient) publicMessage(funcName string, from string, obj interface{}) {
	var arg []byte
	if obj != nil {
		var err error
		arg, err = json.Marshal(obj)
		if err != nil {
			log.Error(err.Error())
		}
	}
	log.Info(string(arg))
	callInfo := m.buildDefaultCallInfo(funcName, from, arg)
	m.rpcClient.Call(*callInfo, m.callback_chan)
}

func (m DefaultMasterClient) startListen(callChan chan mqrpc.CallInfo, selfCallChan chan mqrpc.CallInfo) {
	for {
		select {
		case callInfo := <-callChan:
			from := callInfo.RpcInfo.ReplyTo
			if from != m.Name {
				funcName := callInfo.RpcInfo.Fn
				switch funcName {
				case OnRegister:
					appInfo := &ApplicationInfo{
					}
					err := json.Unmarshal(callInfo.RpcInfo.Args, appInfo)
					if err != nil {
						panic(err)
					}
					m.updateModuleInfos(appInfo)

				case OnUnRegister:
					name := ""
					err := json.Unmarshal(callInfo.RpcInfo.Args, name)
					if err != nil {
						panic(err)
					}

				case UpdateStatus:
					log.Info("update status %s", string(callInfo.RpcInfo.Args))
					apps := []AppStatus{}
					err := json.Unmarshal(callInfo.RpcInfo.Args, &apps)
					if err != nil {
						panic(err)
					}

				}
			}


		case selfCallInfo := <-selfCallChan:
			funcName := selfCallInfo.RpcInfo.Fn
			log.Info(funcName)
			switch funcName {
			case getAppList:
				infos := []*ApplicationInfo{}
				error := json.Unmarshal(selfCallInfo.RpcInfo.Args, &infos)
				if error != nil {
					log.Error(error.Error())
				}
				for _, appInfo := range infos {
					m.updateModuleInfos(appInfo)
				}
			case GetAppStatus:

			}
		}
	}
}

func (m DefaultMasterClient) updateModuleInfos(appInfo *ApplicationInfo) {
	if appInfo.Name != m.Name {
		for moduleName, modules := range appInfo.Modules {
			for _, module := range modules {
				rpcClient, err := defaultrpc.NewRPCClient(m.app, m.Name)
				if err != nil {
					panic(err)
				}

				if module.Rabbitmq != nil {
					//如果远程的rpc存在则创建一个对应的客户端
					rpcClient.NewRabbitmqClient(module.Rabbitmq)
				}

				if module.Redis != nil {
					//如果远程的rpc存在则创建一个对应的客户端
					rpcClient.NewRedisClient(module.Redis)
				}

				serverSession := basemodule.NewServerSession(module.Id, moduleName, module.ByteID, rpcClient)
				moduleInfo := &OtherModuleInfo{
					serverSession: &serverSession,
					appName:       m.Name,
					key:           m.Name + module.Id,
				}
				m.checkToRemoveFromCacheModule(moduleInfo)
				m.lock.Lock()
				m.moduleInfo[moduleName] = append(m.moduleInfo[moduleName], moduleInfo)
				m.moduleInfoByte[module.ByteID] = append(m.moduleInfoByte[module.ByteID], moduleInfo)
				m.lock.Unlock()
			}
		}
	}
}

func (m DefaultMasterClient) checkToRemoveFromCacheModule(moduleInfo *OtherModuleInfo) {
	m.lock.Lock()
	for key, value := range m.moduleInfo {
		for index, inModuleInfo := range value {
			if inModuleInfo.key == moduleInfo.key {
				m.moduleInfo[key] = append(value[:index], value[index+1:]...)
			}
		}
	}
	for key, value := range m.moduleInfoByte {
		for index, inModuleInfo := range value {
			if inModuleInfo.key == moduleInfo.key {
				m.moduleInfoByte[key] = append(value[:index], value[index+1:]...)
			}
		}
	}
	m.lock.Unlock()
}

func (m DefaultMasterClient) removeApplication(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for key, value := range m.moduleInfo {
		for index, inModuleInfo := range value {
			if inModuleInfo.appName == name {
				m.moduleInfo[key] = append(value[:index], value[index+1:]...)
			}
		}
	}

	for key, value := range m.moduleInfo {
		for index, inModuleInfo := range value {
			if inModuleInfo.appName == name {
				m.moduleInfo[key] = append(value[:index], value[index+1:]...)
			}
		}
	}
}

func (m DefaultMasterClient) tick() {
	for {
		time.Sleep(1 * time.Second)
		moduleInfo := make([]ModuleStatus, len(m.moduleManager.GetModules()))
		var allLoad int32
		var allProcessingNumbers int32
		for index, subModule := range m.moduleManager.GetModules() {
			m := subModule.Mi
			if m == nil {
				log.Error("Module is nil")
				return
			}
			rpcModule, b := subModule.Mi.(module.RPCModule)
			if b {
				statistical, _ := rpcModule.GetStatistical()
				var load int32
				load = int32(rpcModule.GetExecuting())
				moduleInfo[index] = ModuleStatus{
					ModuleName:  subModule.Mi.GetType(),
					Load:        load,
					MethodLoads: statistical,
				}
				allLoad += load
				allProcessingNumbers += load
			}
		}

		appStatus := &AppStatus{
			AppName:           m.Name,
			Load:              allLoad,
			ProcessIngNumbers: allProcessingNumbers,
			ModuleStatus:      moduleInfo,
		}
		m.publicMessage(ReportStatus, m.Name, appStatus)
	}

}
