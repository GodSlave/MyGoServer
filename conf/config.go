// Copyright 2014 mqant Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	LenStackBuf = 1024
	Conf        = Config{}
)

func LoadConfig(Path string) {
	// Read config.
	if err := readFileInto(Path); err != nil {
		panic(err)
	}

}

type Config struct {
	Name        string
	Rpc         Rpc
	Module      map[string][]*ModuleSettings
	Master      Master
	Debug       bool
	DB          DB
	PrivateKey  string
	Secret      bool
	OnlineLimit int32
}

type DB struct {
	DBtype string
	SQL    string
	Redis  string
}

type Rpc struct {
	RpcExpired int //远程访问最后期限值 单位秒[默认5秒] 这个值指定了在客户端可以等待服务端多长时间来应答
}

type Rabbitmq struct {
	Uri          string
	Exchange     string
	ExchangeType string
	Queue        string
	BindingKey   string //
	ConsumerTag  string //消费者TAG
}

type Redis struct {
	Uri   string //redis://:[password]@[ip]:[port]/[db]
	Queue string
}

type ModuleSettings struct {
	Id        string
	ByteID    byte
	Host      string
	ProcessID string
	Settings  map[string]interface{}
	Rabbitmq  *Rabbitmq
	Redis     *Redis
}

type SSH struct {
	Host     string
	Port     int
	User     string
	Password string
}

/**
host:port
*/
func (s *SSH) GetSSHHost() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Process struct {
	ProcessID string
	Host      string
	//执行文件
	Execfile string
	//日志文件目录
	//pid.nohup.log
	//pid.access.log
	//pid.error.log
	LogDir string
	//自定义的参数
	Args map[string]interface{}
}

type Master struct {
	ISRealMaster    bool
	RedisUrl        string
	DBConfig        *InfluxDBConfig
	RedisPubSubConf *Redis
	Enable          bool
}

type InfluxDBConfig struct {
	DBName    string
	UserName  string
	Password  string
	Addr      string
	Precision string
	Enable    bool
}

func readFileInto(path string) error {
	var data []byte
	buf := new(bytes.Buffer)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if len(line) > 0 {
				buf.Write(line)
			}
			break
		}
		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
			buf.Write(line)
		}
	}
	data = buf.Bytes()
	return json.Unmarshal(data, &Conf)
}

// If read the file has an error,it will throws a panic.
func fileToStruct(path string, ptr *[]byte) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	*ptr = data
}
