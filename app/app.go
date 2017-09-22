package app

import (
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"hash/crc32"
	"math"
	"os/exec"
	"os"
	"path/filepath"
	"fmt"
	"flag"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/module/modules"
	"os/signal"
	"github.com/GodSlave/MyGoServer/rpc/base"
	"github.com/GodSlave/MyGoServer/rpc"
	"strings"
	"github.com/GodSlave/MyGoServer/db"
	"github.com/go-xorm/xorm"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/garyburd/redigo/redis"
	"github.com/GodSlave/MyGoServer/utils"
	"encoding/json"
)

func NewApp() module.App {
	newApp := new(DefaultApp)
	newApp.routes = map[string]func(app module.App, Type string, hash string) module.ServerSession{}
	newApp.byteRoutes = map[byte]func(app module.App, Type byte, hash string) module.ServerSession{}
	newApp.byteServerList = map[byte]module.ServerSession{}
	newApp.serverList = map[string]module.ServerSession{}
	newApp.defaultRoutes = func(app module.App, Type string, hash string) module.ServerSession {
		servers := app.GetServersByType(Type)
		if len(servers) == 0 {
			log.Error("no smodule find %s", Type)
			return nil
		}
		index := int(math.Abs(float64(crc32.ChecksumIEEE([]byte(hash))))) % len(servers)
		return servers[index]
	}

	newApp.byteDefaultRoutes = func(app module.App, Type byte, hash string) module.ServerSession {
		servers := app.GetServersByByteType(Type)
		if len(servers) == 0 {
			log.Error("no module find %v", Type)
			return nil
		}
		index := int(math.Abs(float64(crc32.ChecksumIEEE([]byte(hash))))) % len(servers)
		return servers[index]
	}

	newApp.rpcserializes = map[string]module.RPCSerialize{}
	newApp.version = "0.0.1"
	newApp.users = map[string]*base.BaseUser{}
	return newApp
}

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
	users             map[string]*base.BaseUser
	psc               *redis.PubSubConn
}

func (app *DefaultApp) Run(mods ...module.Module) error {
	file, _ := exec.LookPath(os.Args[0])
	ApplicationPath, _ := filepath.Abs(file)
	ApplicationDir, _ := filepath.Split(ApplicationPath)
	defaultPath := fmt.Sprintf("%sconf/server.json", ApplicationDir)
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
	log.Info("mqant %v starting up", app.version)

	log.Info("start connect DB %v", conf.Conf.DB.SQL)
	//sql
	sql := db.BaseSql{
	}
	sql.Url = conf.Conf.DB.SQL
	sql.InitDB()
	sql.CheckMigrate()
	app.Engine = sql.Engine
	defer app.Engine.Close()

	url := app.GetSettings().DB.Redis
	app.redisPool = utils.GetRedisFactory().GetPool(url)
	defer app.redisPool.Close()

	psc := redis.PubSubConn{Conn: app.redisPool.Get()}
	psc.Subscribe(app.GetSettings().DB.Redis_Queue)
	app.psc = &psc
	go app.handAppMessage()

	log.Info("start register module %v", conf.Conf.DB.SQL)

	manager := basemodule.NewModuleManager()
	manager.RegisterRunMod(modules.TimerModule())

	// module
	for i := 0; i < len(mods); i++ {
		manager.Register(mods[i])
	}
	app.OnInit(app.settings)
	manager.Init(app, *ProcessID)
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	manager.Destroy()
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

func (app *DefaultApp) RpcInvokeArgs(module module.RPCModule, moduleType string, _func string, sessionId string, args []byte) (result interface{}, err *base.ErrorCode) {
	server, e := app.GetRouteServers(moduleType, module.GetServerId())
	if e != nil {
		err = base.NewError(404, e.Error())
		return
	}
	return server.CallArgs(_func, sessionId, args)
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

func (app *DefaultApp) VerifyUserID(sessionId string) (userID string) {
	c := app.redisPool.Get()
	reply, err := redis.Bool(c.Do("EXISTS", base.TOKEN_PERFIX+sessionId))
	if err == nil && reply {

		userID, _ = redis.String(c.Do("GET", base.TOKEN_PERFIX+sessionId))
		return
	}
	reply, err = redis.Bool(c.Do("EXISTS", base.SESSION_PERFIX+sessionId))
	if err == nil && reply {
		userID, _ = redis.String(c.Do("GET", base.SESSION_PERFIX+sessionId))
		return
	}
	return ""
}

func (app *DefaultApp) VerifyUser(sessionId string) (user *base.BaseUser) {

	var exit bool
	user, exit = app.users[sessionId]
	if exit {
		return
	}
	uid := app.VerifyUserID(sessionId)
	if uid == "" {
		return nil
	}
	user = &base.BaseUser{
		UserID: uid,
	}
	result, err := app.Engine.Get(user)
	if err != nil {
		panic(err)
	}
	if !result {
		return nil
	}
	app.users[sessionId] = user
	return
}

func (app *DefaultApp) OnUserLogin(sessionId string) {
	app.VerifyUser(sessionId)
}

type ProcessMessageContent struct {
	Method string
	Body   interface{}
}

const LOGOUT_MESSAGE = "OnUserLogOut"

func (app *DefaultApp) OnUserLogOut(sessionId string) {
	delete(app.users, sessionId)
	redis := app.redisPool.Get()
	redis.Do("DEL", base.SESSION_PERFIX+sessionId)
	message := ProcessMessageContent{
		Method: LOGOUT_MESSAGE,
		Body:   sessionId,
	}
	msgData, err := json.Marshal(message)
	if err != nil {
		log.Error(err.Error())
		return
	}
	redis.Do("PUBLISH", app.settings.DB.Redis_Queue, msgData)
}

func (app *DefaultApp) handAppMessage() {

	for {
		switch v := app.psc.Receive().(type) {
		case redis.Message:
			appMsg := &ProcessMessageContent{}
			err := json.Unmarshal(v.Data, appMsg)
			if err != nil {
				log.Error(err.Error())
			}
			switch appMsg.Method {
			case LOGOUT_MESSAGE:
				sessID := appMsg.Body.(string)
				delete(app.users, sessID)
			}
		case redis.PMessage:
		case redis.Subscription:
			log.Info("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Error("on_request_handle", v.Error())
			return
		default:

		}

	}
}
