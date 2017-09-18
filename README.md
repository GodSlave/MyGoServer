Go 语言的游戏引擎

在（https://github.com/liangdas/mqant）mquant的基础上修改了一部分代码

1. rpc部分没有动
2. 替换了mqtt模块
3. 替换了rpc方法调用时 map[string]interface 为具体的struct   
4  增加了基本的用户模块


protoBuf generate sample
protoc --gofast_out=. -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf/ user.proto
