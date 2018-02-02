## 概述

Go游戏服务器  使用了 protoBuf  MQTT 等协议
在mqant基础上进行了部分改进

1. 支持标准mqtt协议，支持标准mqtt客户端连接,通讯
2. 内置redis,MySQL 数据存储支持
3. 内置基本的用户模块  


游戏业务服务使用MQTT与客户端通讯，
MQTT服务端使用surgeMQ 并为了安全性和包大小考虑做了一部分修改
1. 限制clientID只能32-64位字符串
2. 用户先注册名称为s的topic用来授权和获取密钥
3. 获取密钥成功后可以注册f，i 的topic， 
4. 限制每个链接的clientID不能相同
5. s f i topic 全部采用qos1 消息级别
6. 因为目前内置mqtt还不支持SSL ，所以登录操作在获取密钥后进行


f  topic 表示通讯内容使用protobuf 格式  
i  topic 表示通讯内容使用json格式
p  topic 表示服务器发来的推送消息 json格式
q  topic 表示服务器发来的推送消息 protobuf格式

注意： protoBuf格式的推送会在p和q通道同时推送，请根据需要注册

注册s topic成功后  服务器端会返回一段32位字符串，使用 服务器端字符串和客户端ID进行字符串拼接 并交换第32位和第6位，交换第48位和第31位后进行MD5取值，并取md5结果16位byte作为AES加密密钥 同时作为iv和ciper， 获取加密密钥后可以注册s端口
获取AES密钥成功后，后续 s，f topic 通讯统一使用该密钥对payload内容加密  

推荐使用idea的goland IDE开发

目前已经实现了基本的user模块，对照上手非常简单，欢迎使用，交流



## topic 定义
  topic s i f 开头的频道不能用于普通MQTT通信，需注意

## message 定义
### 请求json格式的message 为
    type MsgFormat struct {
    	Module string      `json:"module"`
    	Func   string      `json:"func"`
	Params interface{} `json:"params"`
    }

### protobuf 格式为
byte array [----...]   
byte1 对应 Module编号  
byte2 对应 Func编号   
剩余数据对应 参数的protobuf的数据
### 反馈数据
    
    message allResponse {

        //结果
        bytes result = 1;
        //错误结果 如果为nil表示请求正确
        string msg = 2;
        //错误代码 200 means right
        int32 state = 3;

    }

如果请求正确， state 为 0    
如果请求错误， state 从1 开始为该方法特定错误，一般需要单独处理， 从ff倒排表示 通用错误

## 方法定义
### 方法需要在模块初始化时注册  
 `m.GetServer().RegisterGO("Login", 1, m.Login)`
### 方法的参数定义  

 `func (m *ModuleUser) Login(SessionId string, form *UserLoginRequest) (result *UserLoginResponse, err *base.ErrorCode)`  

###  参数 
 方法支持两个参数和一个参数，
 参数1. 可以为string格式 表示 userID  如果用户没有登录或者token不对 会传递 sessionID  
       也可以为base.BaseUser 调用模块会自动注入user实体  
       如果不需要参数1 可以不写  
 参数2. 直接为对应请求bean的实体 方便使用，如果不需求直接刘空  
### 返回  
 result 表示要返回的实体bean  
 err 表示处理过程中需要的返回的错误，对客户端来说如果code不正确 不需要处理返回的实体bean  
 这两个返回都可以为 空  
      

## 如何开始
 1. 配置conf目录下server.json文件  
 2. go build main.go  & ./main  
 3. 进入test目录 找到 testJson 方法 开始执行

## 模块列表
| module | module |
| -------- | -------- |
|gate|1|
|user|2|
|push|3|
|push|4|
... 保留编号到10
请用户使用的时候从11开始使用

## 负载均衡  
参考work目录下的work1-work4   
work1 为主进程，实现子进程管理   对应main1.go
work2 为网关模块，负责用户接入   对应main2.go
work3，work4分别为用户模块，     对应main3.go

运行test目录中的mqtt.test 的TestJson 方法 可以看到work3和work4 同时在处理用户的注册请求  

服务状态监控平台正在增加中...








