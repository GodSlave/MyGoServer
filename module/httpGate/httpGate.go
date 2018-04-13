package httpGate

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"net/http"
	"github.com/GodSlave/MyGoServer/log"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"encoding/json"
)

type HttpHandler struct {
	http.Handler
	httpGate             *HttpGate
	errSerializationFail []byte
	errModuleNotFound    []byte
}

type Response struct {
	Status int32       `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func (handler *HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Info(fmt.Sprintf("%v", request))
	var module string
	var method string
	var param []byte
	var session string

	path := request.URL.Path
	moduleInfos := strings.Split(path, "/")
	if len(moduleInfos) > 2 {
		module = moduleInfos[1]
		method = moduleInfos[2]
	} else {
		writer.WriteHeader(500)
		writer.Write([]byte("url format error"))
		return
	}
	session = request.Header.Get("Session")
	if session == "" {
		session = request.Header.Get("cookie")
	}

	if session == "" {
		session = uuid.SafeString(32)
		writer.Header().Set("Session", session)
		writer.Header().Set("cookie", session)
	}

	if request.ContentLength > 0 {
		defer request.Body.Close()
		body, err := ioutil.ReadAll(request.Body)
		param = body
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(string(body))
	}

	callSession, error := handler.httpGate.app.GetRouteServers(module, "")
	if error == nil && callSession != nil {
		result, errCode := callSession.CallArgs(method, session, param)
		restlutformat := `{"status":%d,data:%s,msg:"%s"}`
		result1 := fmt.Sprintf(restlutformat, errCode.ErrorCode, (string(result)), errCode.Desc)
		writer.WriteHeader(200)
		writer.Write([]byte(result1))
	} else {
		writer.WriteHeader(500)
		writer.Write(handler.errModuleNotFound)
	}
}

type HttpGate struct {
	basemodule.BaseModule
	app         module.App
	listenUrl   string
	httpHandler *HttpHandler
}

var Module = func() module.Module {
	newModule := new(HttpGate)
	return newModule
}

func (this *HttpGate) OnInit(app module.App, settings *conf.ModuleSettings) {
	this.BaseModule.OnInit(this, app, settings) //这是必须的
	this.app = app
	this.listenUrl = (settings.Settings["listenUrl"]).(string)
	this.httpHandler = &HttpHandler{
		httpGate: this,
	}

	errResponse := &Response{
		Status: 500,
		Msg:    "serialization fail",
	}

	this.httpHandler.errSerializationFail, _ = json.Marshal(errResponse)
	errResponse.Msg = "Module Not Found"
	this.httpHandler.errModuleNotFound, _ = json.Marshal(errResponse)

}

func (this *HttpGate) GetType() string {
	return "HttpGate"
}

func (this *HttpGate) Version() string {
	return "1.0.0"
}

func (this *HttpGate) Run(closeSig chan bool) {
	http.ListenAndServe(this.listenUrl, this.httpHandler)
}

func (this *HttpGate) OnDestroy() {

}
