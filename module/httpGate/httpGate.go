package httpGate

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"net/http"
	"github.com/GodSlave/MyGoServer/log"
	"fmt"
	"io/ioutil"
)

type HttpHandler struct {
	http.Handler
	httpGate *HttpGate
}

func (handler *HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var module string
	var method string
	var param []byte




	log.Info(request.Method)
	log.Info(request.URL.Path)
	log.Info(request.Header.Get("Session"))
	log.Info(fmt.Sprintf("%v", request))
	log.Info(fmt.Sprintf("%v", request.Form))
	log.Info(fmt.Sprintf("%v", request.URL.RawQuery))
	log.Info(fmt.Sprintf("%v", request.PostForm))
	if request.ContentLength > 0 {
		defer request.Body.Close()
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(string(body))
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
