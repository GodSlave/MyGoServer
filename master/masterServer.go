package master

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/GodSlave/MyGoServer/utils"
	"github.com/GodSlave/MyGoServer/db"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/rpc/base"
	"github.com/GodSlave/MyGoServer/rpc"
	"github.com/GodSlave/MyGoServer/rpc/pb"
	"strconv"
	"time"
	"encoding/json"
)

type DefaultMasterServer struct {
	Master
	Engine        *xorm.Engine
	redisPool     *redis.Pool
	rpcClient     *defaultrpc.RedisClient
	rpcServer     *defaultrpc.RedisServer
	infos         map[string]ApplicationInfo
	status        map[string]AppStatus
	appRpcClient  map[string]*defaultrpc.RedisClient
	callback_chan chan rpcpb.ResultInfo
	call_chan     chan mqrpc.CallInfo
	serverId      string
	callId        int
	versionCode   int32
	masterConfig  conf.Master
	Modules       map[string][]OtherModuleInfo
}

func NewMaster(serverId string, masterConf conf.Master) Master {
	master := new(DefaultMasterServer)
	master.masterConfig = masterConf
	master.infos = map[string]ApplicationInfo{}
	master.status = map[string]AppStatus{}
	master.appRpcClient = map[string]*defaultrpc.RedisClient{}
	redisUrl := masterConf.RedisUrl
	master.redisPool = utils.GetRedisFactory().GetPool(redisUrl)
	master.serverId = serverId
	sql := db.BaseSql{
	}
	SQlUrl := masterConf.SQlUrl
	sql.Url = SQlUrl
	log.Info(sql.Url)
	sql.InitDB()
	sql.CheckMigrate()
	master.Engine = sql.Engine

	master.rpcClient, _ = defaultrpc.NewRedisClient(masterConf.RedisPubSubConf)
	master.call_chan = make(chan mqrpc.CallInfo, 10)
	master.rpcServer, _ = defaultrpc.NewRedisServer(masterConf.RedisPubSubConf, master.call_chan)
	master.callback_chan = make(chan rpcpb.ResultInfo, 10)

	go master.checkReceiverMessage()
	go master.tickServerStatus()

	return master
}

func (m *DefaultMasterServer) Register(info ApplicationInfo) {
	m.infos[info.Name] = info
	redisConf := &conf.Redis{
		Uri:   m.masterConfig.RedisPubSubConf.Uri,
		Queue: info.Name,
	}
	rpcClient, _ := defaultrpc.NewRedisClient(redisConf)
	m.appRpcClient[info.Name] = rpcClient
	m.versionCode += 1
	m.publicMessage(OnRegister, info.Name, info, m.rpcClient)
	m.GetAppList(info.Name)
}

func (m *DefaultMasterServer) UnRegister(appName string) {
	delete(m.infos, appName)
	m.versionCode += 1
	m.publicMessage(OnUnRegister, appName, appName, m.rpcClient)
}

func (m *DefaultMasterServer) GetAppStatus(name string) []AppStatus {
	if val, ok := m.appRpcClient[name]; ok {
		tempStatus := make([]AppStatus, len(m.status))
		index := 0
		for _, value := range m.status {
			value.ModuleStatus = nil
			tempStatus[index] = value
			index ++
		}
		m.publicMessage(GetAppStatus, MasterStr, tempStatus, val)
	}
	return nil
}

func (m *DefaultMasterServer) GetAppList(name string) []ApplicationInfo {
	if val, ok := m.appRpcClient[name]; ok {
		tempApp := make([]ApplicationInfo, len(m.infos))
		index := 0
		for _, value := range m.infos {
			tempApp[index] = value
			index ++
		}
		m.publicMessage(getAppList, MasterStr, tempApp, val)
	}
	return nil
}

func (m *DefaultMasterServer) ReportStatus(status AppStatus) {
	m.status[status.AppName] = status
}

func (m *DefaultMasterServer) buildDefaultCallInfo(functionName string, from string, args []byte) *mqrpc.CallInfo {
	m.callId += 1
	callInfo := mqrpc.CallInfo{
		RpcInfo: rpcpb.RPCInfo{
			Cid:       m.serverId + strconv.Itoa(m.callId),
			Fn:        functionName,
			Args:      args,
			Reply:     false,
			SessionId: from,
			ByteFn:    m.versionCode,
		},
	}

	return &callInfo;
}

func (m *DefaultMasterServer) publicMessage(funcName string, from string, obj interface{}, client *defaultrpc.RedisClient) {
	var arg []byte
	if obj != nil {
		var err error
		arg, err = json.Marshal(obj)
		if err != nil {
			log.Error(err.Error())
		}
	}
	log.Info("published %s ", string(arg))
	callInfo := m.buildDefaultCallInfo(funcName, from, arg)
	client.Call(*callInfo, m.callback_chan)
}

func (m *DefaultMasterServer) checkReceiverMessage() {
	for {
		select {
		case callInfo := <-m.call_chan:
			functionName := callInfo.RpcInfo.Fn
			switch functionName {
			case Register:
				info := &ApplicationInfo{}
				json.Unmarshal(callInfo.RpcInfo.Args, info)
				m.Register(*info)
			case UnRegister:
				appName := string(callInfo.RpcInfo.Args)
				m.UnRegister(appName)
			case GetAppStatus:
				appName := string(callInfo.RpcInfo.Args)
				m.GetAppStatus(appName)
			case getAppList:
				appName := string(callInfo.RpcInfo.Args)
				m.GetAppList(appName)
			case ReportStatus:
				appStatus := &AppStatus{}
				err := json.Unmarshal(callInfo.RpcInfo.Args, appStatus)
				if err != nil {
					log.Error(err.Error())
					return
				}
				m.status[appStatus.AppName] = *appStatus
			}
			//case resultInfo := <-m.callback_chan:
			//log.Error("%v", resultInfo.Cid)
		}
	}
}
func (m *DefaultMasterServer) tickServerStatus() {
	for {
		time.Sleep(1 * time.Second)
		tempStatus := make([]AppStatus, len(m.status))
		index := 0
		for _, value := range m.status {
			value.ModuleStatus = nil
			tempStatus[index] = value
			index ++
		}
		m.publicMessage(UpdateStatus, MasterStr, tempStatus, m.rpcClient)
		//m.publicMessage(SyncVersionCode, MasterStr, m.versionCode, m.rpcClient) //keep client  & master  in same status
	}
}
