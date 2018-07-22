package app

import (
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"os/exec"
	"os"
	"path/filepath"
	"fmt"
	"flag"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/rpc/base"
	"github.com/GodSlave/MyGoServer/rpc"
	"strings"
	"github.com/GodSlave/MyGoServer/db"
	"github.com/go-xorm/xorm"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/utils"
	"time"
	"github.com/GodSlave/MyGoServer/master"
	"os/signal"
	"math"
	"hash/crc32"
	"github.com/GodSlave/MyGoServer/base"
)

type DefaultApp struct {
	module.App
	version           string
	serverList        map[string]module.ServerSession
	byteServerList    map[byte]module.ServerSession
	settings          conf.Config
	routes            map[string]func(app module.App, Type string, hash string) module.ServerSession
	byteRoutes        map[byte]func(app module.App, Type byte, hash string) module.ServerSession
	defaultRoutes     func(app module.App, Type string, hash string) module.ServerSession
	byteDefaultRoutes func(app module.App, Type byte, hash string) module.ServerSession
	rpcserializes     map[string]module.RPCSerialize
	Engine            *xorm.Engine
	redisPool         *redis.Pool
	psc               *redis.PubSubConn
	userManager       module.UserManager
	initDownCallback  module.OnInitDownCallBack
	masterServer      master.Master
	masterClient      master.MasterClient
	moduleManger      *basemodule.ModuleManager
}

func NewApp() module.App {
	newApp := new(DefaultApp)
	newApp.routes = map[string]func(app module.App, Type string, hash string) module.ServerSession{}
	newApp.byteRoutes = map[byte]func(app module.App, Type byte, hash string) module.ServerSession{}
	newApp.byteServerList = map[byte]module.ServerSession{}
	newApp.serverList = map[string]module.ServerSession{}
	newApp.defaultRoutes = func(app module.App, Type string, hash string) module.ServerSession {
		if newApp.masterClient != nil {
			serverSession := newApp.masterClient.GetModule(Type)
			if serverSession != nil {
				return *serverSession
			}
		} else {
			servers := app.GetServersByType(Type)
			if len(servers) == 0 {
				log.Error("no smodule find %s", Type)
				return nil
			}
			index := int(math.Abs(float64(crc32.ChecksumIEEE([]byte(hash))))) % len(servers)
			return servers[index]
		}
		return nil
	}

	newApp.byteDefaultRoutes = func(app module.App, Type byte, hash string) module.ServerSession {

		if newApp.masterClient != nil {
			return *newApp.masterClient.GetModuleByByte(Type)
		} else {
			servers := app.GetServersByByteType(Type)
			if len(servers) == 0 {
				log.Error("no module find %v", Type)
				return nil
			}
			index := int(math.Abs(float64(crc32.ChecksumIEEE([]byte(hash))))) % len(servers)
			return servers[index]
		}
		return nil
	}

	newApp.rpcserializes = map[string]module.RPCSerialize{}
	newApp.version = "0.0.1"
	return newApp
}

func (app *DefaultApp) Run(mods ...module.Module) error {
	file, _ := exec.LookPath(os.Args[0])
	ApplicationPath, _ := filepath.Abs(file)
	ApplicationDir, _ := filepath.Split(ApplicationPath)
	defaultPath := fmt.Sprintf("%sconf"+string(filepath.Separator)+"server.json", ApplicationDir)
	confPath := flag.String("conf", defaultPath, "Server configuration file path")
	ProcessID := flag.String("pid", "development", "Server ProcessID?")
	Logdir := flag.String("log", fmt.Sprintf("%slogs", ApplicationDir), "Log file directory?")
	flag.Parse() //解析输入的参数
	f, err := os.Open(*confPath)
	if err != nil {
		panic(err)
	}
	_, err = os.Open(*Logdir)
	if err != nil {
		//文件不存在
		err := os.Mkdir(*Logdir, os.ModePerm) //
		if err != nil {
			fmt.Println(err)
		}
	}
	log.Info("Server configuration file path [%s]", *confPath)
	conf.LoadConfig(f.Name()) //加载配置文件
	app.Configure(conf.Conf)  //配置信息
	log.Init(conf.Conf.Debug, *ProcessID, *Logdir)
	log.Info("server %v starting up at %v", app.version, time.Now().Unix())
	log.Debug("start connect DB %v", conf.Conf.DB.SQL)
	app.userManager = InitUserManager(app, conf.Conf.OnlineLimit)

	//sql
	sql := db.BaseSql{
	}
	sql.Url = conf.Conf.DB.SQL
	log.Info(sql.Url)
	sql.InitDB()
	sql.CheckMigrate()
	app.Engine = sql.Engine
	defer app.Engine.Close()

	url := app.GetSettings().DB.Redis
	app.redisPool = utils.GetRedisFactory().GetPool(url)
	defer app.redisPool.Close()

	// module
	log.Info("start register module %v", conf.Conf.DB.SQL)
	app.moduleManger = basemodule.NewModuleManager()

	for i := 0; i < len(mods); i++ {
		app.moduleManger.Register(mods[i])
	}

	if conf.Conf.Master.ISRealMaster && conf.Conf.Master.Enable {
		app.masterServer = master.NewMaster(conf.Conf.Name, conf.Conf.Master)
	}
	if conf.Conf.Master.Enable {
		app.masterClient = master.NewMasterClient(conf.Conf.Master, conf.Conf.Name, app, app.moduleManger)
	}

	app.OnInit(app.settings)
	app.moduleManger.Init(app, *ProcessID)

	if app.initDownCallback != nil {
		app.initDownCallback(app)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	app.moduleManger.Destroy()
	app.OnDestroy()

	log.Info("mqant closing down (signal: %v)", sig)
	return nil
}

func (app *DefaultApp) Route(moduleType string, fn func(app module.App, Type string, hash string) module.ServerSession) error {
	app.routes[moduleType] = fn
	return nil
}
func (app *DefaultApp) getRoute(moduleType string) func(app module.App, Type string, hash string) module.ServerSession {
	fn := app.routes[moduleType]
	if fn == nil {
		//如果没有设置的路由,则使用默认的
		return app.defaultRoutes
	}
	return fn
}

func (app *DefaultApp) getByteRoute(moduleType byte) func(app module.App, Type byte, hash string) module.ServerSession {
	fn := app.byteRoutes[moduleType]
	if fn == nil {
		//如果没有设置的路由,则使用默认的
		return app.byteDefaultRoutes
	}
	return fn
}

func (app *DefaultApp) AddRPCSerialize(name string, Interface module.RPCSerialize) error {
	if _, ok := app.rpcserializes[name]; ok {
		return fmt.Errorf("The name(%s) has been occupied", name)
	}
	app.rpcserializes[name] = Interface
	return nil
}

func (app *DefaultApp) GetRPCSerialize() (map[string]module.RPCSerialize) {
	return app.rpcserializes
}

func (app *DefaultApp) Configure(settings conf.Config) error {
	app.settings = settings
	return nil
}

/**

 */
func (app *DefaultApp) OnInit(settings conf.Config) error {
	app.serverList = make(map[string]module.ServerSession)
	for Type, ModuleInfos := range settings.Module {
		for _, moduel := range ModuleInfos {
			m := app.serverList[moduel.Id]
			if m != nil {
				//如果Id已经存在,说明有两个相同Id的模块,这种情况不能被允许,这里就直接抛异常 强制崩溃以免以后调试找不到问题
				panic(fmt.Sprintf("ServerId (%s) Type (%s) of the modules already exist Can not be reused ServerId (%s) Type (%s)", m.GetId(), m.GetType(), moduel.Id, Type))
			}
			client, err := defaultrpc.NewRPCClient(app, moduel.Id)
			if err != nil {
				continue
			}
			if moduel.Rabbitmq != nil {
				//如果远程的rpc存在则创建一个对应的客户端
				client.NewRabbitmqClient(moduel.Rabbitmq)
			}
			if moduel.Redis != nil {
				//如果远程的rpc存在则创建一个对应的客户端
				client.NewRedisClient(moduel.Redis)
			}
			session := basemodule.NewServerSession(moduel.Id, Type, moduel.ByteID, client)
			app.serverList[moduel.Id] = session
			app.byteServerList[moduel.ByteID] = session
			log.Info("RPCClient create success type(%s) id(%s)", Type, moduel.Id)
		}
	}
	return nil
}

func (app *DefaultApp) OnDestroy() error {
	for id, session := range app.serverList {
		err := session.GetRpc().Done()
		if err != nil {
			log.Warning("RPCClient close fail type(%s) id(%s)", session.GetType(), id)
		} else {
			log.Info("RPCClient close success type(%s) id(%s)", session.GetType(), id)
		}
	}
	return nil
}

func (app *DefaultApp) RegisterLocalClient(serverId string, server mqrpc.RPCServer) error {
	if session, ok := app.serverList[serverId]; ok {
		return session.GetRpc().NewLocalClient(server)
	} else {
		return fmt.Errorf("Server(%s) Not Found", serverId)
	}
	return nil
}

func (app *DefaultApp) GetServersById(serverId string) (module.ServerSession, error) {
	if session, ok := app.serverList[serverId]; ok {
		return session, nil
	} else {
		return nil, fmt.Errorf("Server(%s) Not Found", serverId)
	}
}

func (app *DefaultApp) GetServersByType(Type string) []module.ServerSession {
	sessions := make([]module.ServerSession, 0)
	for _, session := range app.serverList {
		if session.GetType() == Type {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (app *DefaultApp) GetServersByByteType(Type byte) []module.ServerSession {
	sessions := make([]module.ServerSession, 0)
	for _, session := range app.serverList {
		if session.GetByteType() == Type {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (app *DefaultApp) GetRouteServers(filter string, hash string) (s module.ServerSession, err error) {
	sl := strings.Split(filter, "@")
	if len(sl) == 2 {
		moduleID := sl[1]
		if moduleID != "" {
			return app.GetServersById(moduleID)
		}
	}
	moduleType := sl[0]
	route := app.getRoute(moduleType)
	if route == nil {
		return nil, base.ErrFrozen
	}
	s = route(app, moduleType, hash)
	if s == nil {
		log.Error("Server(type : %s) Not Found", moduleType)
	}
	return
}

func (app *DefaultApp) GetByteRouteServers(filter byte, hash string) (s module.ServerSession, err error) {
	route := app.getByteRoute(filter)
	s = route(app, filter, hash)
	if s == nil {
		log.Error("Server(type : %x) Not Found", filter)
	}
	return
}

func (app *DefaultApp) GetSettings() conf.Config {
	return app.settings
}

func (app *DefaultApp) RpcInvokeArgs(module module.RPCModule, moduleType string, _func string, sessionId string, args []byte) (result []byte, err *base.ErrorCode) {
	server, e := app.GetRouteServers(moduleType, module.GetServerId())
	if e != nil {
		err = base.NewError(404, e.Error())
		return
	}
	return server.CallArgs(_func, sessionId, args)
}

func (app *DefaultApp) RpcAllInvokeArgs(module module.RPCModule, moduleType string, _func string, sessionId string, args []byte) (result [][]byte, err []*base.ErrorCode) {
	servers := app.GetServersByType(moduleType)

	err = []*base.ErrorCode{}
	for _, server := range servers {
		resultItem, errItem := server.CallArgs(_func, sessionId, args)
		result = append(result, resultItem)
		err = append(err, errItem)
	}
	return
}

func (app *DefaultApp) RpcInvokeNRArgs(module module.RPCModule, moduleType string, _func string, sessionId string, args []byte) (err error) {
	server, err := app.GetRouteServers(moduleType, module.GetServerId())
	if err != nil {
		return
	}
	return server.CallNRArgs(_func, sessionId, args)
}

func (app *DefaultApp) GetSqlEngine() *xorm.Engine {
	return app.Engine
}

func (app *DefaultApp) GetRedis() *redis.Pool {
	return app.redisPool
}

func (app *DefaultApp) GetUserManager() module.UserManager {
	return app.userManager
}

func (app *DefaultApp) SetInitDownCallBack(callBack module.OnInitDownCallBack) {
	app.initDownCallback = callBack
}

func (app *DefaultApp) GetModuleManager() *basemodule.ModuleManager {
	return app.moduleManger
}
func (app *DefaultApp) GetConfig() conf.Config{
	return  app.settings
}