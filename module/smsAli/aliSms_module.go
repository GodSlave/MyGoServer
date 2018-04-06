package smsAli

import (
	"github.com/GodSlave/MyGoServer/module/base"
	"github.com/GodSlave/MyGoServer/conf"
	"github.com/gwpp/alidayu-go"
	"fmt"
	"github.com/GodSlave/MyGoServer/base"
	"github.com/GodSlave/MyGoServer/module"
	"github.com/GiterLab/aliyun-sms-go-sdk/dysms"
	"github.com/GodSlave/MyGoServer/utils/uuid"
	"github.com/GodSlave/MyGoServer/log"
)

var Module = func() module.Module {
	Timer := new(SmsAli)
	return Timer
}

type SmsAli struct {
	basemodule.BaseModule
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
	this.BaseModule.OnInit(this, app, settings)
	this.app = app
	this.aKey = settings.Settings["aKey"].(string)
	this.secreteAKey = settings.Settings["secreteAKey"].(string)
	this.smsTemplateName = settings.Settings["smsTemplateName"].(string)
	this.smsSignName = settings.Settings["smsSignName"].(string)
	this.smsClient = alidayu.NewTopClient(this.aKey, this.secreteAKey)
	this.GetServer().RegisterGO("sendSmsCode", 1, this.sendSmsCode)
}

func (this *SmsAli) sendSmsCode(form *SendSms_Request) (*SendSms_Response, *base.ErrorCode) {
	dysms.HTTPDebugEnable = true
	dysms.SetACLClient(this.aKey, this.secreteAKey) // dysms.New(ACCESSID, ACCESSKEY)

	// send to one person
	respSendSms, err := dysms.SendSms(uuid.SafeString(16), form.PhoneNumber, this.smsSignName, this.smsTemplateName, fmt.Sprintf(`{"code":%s}`, form.VerifyCode)).DoActionWithException()
	if err != nil {
		log.Info("send sms failed  %v %v", err, respSendSms.Error())
		return nil, base.ErrSMSSendFail
	}
	log.Info("send sms succeed %v", respSendSms.GetRequestID())
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
