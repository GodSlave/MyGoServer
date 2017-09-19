##概述

游戏服务器基于mqant 进行了部分改进
1. 支持标准mqtt协议，支持标准mqtt客户端连接
2. 内置redis,MySQL 数据存储
3. 内置基本的用户模块




游戏业务服务使用MQTT与客户端通讯，
MQTT服务端使用surgeMQ 并为了安全性和包大小考虑做了一部分修改
1. 限制clientID只能32-64位字符串
2. 用户先注册名称为s的topic用来授权和获取密钥
3. 获取密钥成功后可以注册f，i 的topic， 
4. 限制每个链接的clientID不能相同
5. s f i topic 全部采用qos1 消息级别

f  topic 表示通讯内容使用protobuf 格式  
i  topic 表示通讯内容使用json格式  


注册s topic成功后  服务器端会返回一段32位字符串，使用 服务器端字符串和客户端ID进行字符串拼接 并交换第32位和第6位，交换第48位和第31位后进行MD5取值，并取md5结果16位byte作为AES加密密钥 同时作为iv和ciper， 获取加密密钥后可以注册s端口
获取AES密钥成功后，后续 s，f topic 通讯统一使用该密钥对payload内容加密  

推荐使用idea的goland IDE开发

##topic 定义


##message 定义



| module | module | 
| -------- | -------- |
|gate|1|
|[user](user)|2|
|token|3|