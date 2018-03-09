package smsAli

import (
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/gwpp/alidayu-go"
	"github.com/gwpp/alidayu-go/request"
	"fmt"
	"github.com/GodSlave/MyGoServer/log"
	"github.com/GodSlave/MyGoServer/base"
)

var Module = func() module.Module {
	Timer := new(SmsAli)
	return Timer
}

type SmsAli struct {
	module.Module
	app             module.App
	aKey            string
	secreteAKey     string
	smsTemplateName string
	smsSignName     string
	smsClient       *alidayu.TopClient
}

func (m *SmsAli) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "SMS"
}

func (this *SmsAli) OnInit(app module.App, settings *conf.ModuleSettings) {
	this.app = app
	this.aKey = settings.Settings["aKey"].(string)
	this.secreteAKey = settings.Settings["secreteAKey"].(string)

	this.smsClient = alidayu.NewTopClient(this.aKey, this.secreteAKey)
}

func (this *SmsAli) sendSmsCode(form *SendSms_Request) (*SendSms_Response, *base.ErrorCode) {
	req := request.NewAlibabaAliqinFcSmsNumSendRequest()
	req.SmsFreeSignName = this.smsSignName
	req.RecNum = form.PhoneNumber
	req.SmsTemplateCode = this.smsTemplateName
	req.SmsParam = fmt.Sprintf(`{"code":%s}`, form.VerifyCode)
	response, err := this.smsClient.Execute(req)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info(fmt.Sprintf("v%", response))
	return nil, base.ErrNil
}

func (this *SmsAli) Run(closeSig chan bool) {
}

func (this *SmsAli) OnDestroy() {
}

func (this *SmsAli) GetApp() (module.App) {
	return this.app
}

func (this *SmsAli) Version() string {
	return "1.0.0"
}
