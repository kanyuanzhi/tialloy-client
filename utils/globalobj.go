package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type GlobalObj struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`

	ServerHost string `json:"server_host,omitempty"`
	ServerPort int    `json:"host,omitempty"`

	TcpMaxPacketSize    uint32 `json:"tcp_max_packet_size,omitempty"`
	TcpWorkerPoolSize   uint32 `json:"tcp_worker_pool_size,omitempty"`
	TcpMaxWorkerTaskLen uint32 `json:"tcp_max_worker_task_len,omitempty"`

	ReconnectInterval time.Duration `json:"reconnect_interval,omitempty"`

	LogMode bool `json:"log_mode,omitempty"`
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/tialloy_client.json")
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err.Error())
	}

	GlobalLog = logrus.New()
	GlobalLog.SetReportCaller(GlobalObject.LogMode)
	if GlobalObject.LogMode == true {
		GlobalLog.SetLevel(logrus.TraceLevel)
	} else {
		GlobalLog.SetLevel(logrus.InfoLevel)
	}
	GlobalLog.SetFormatter(&customFormatter{})
}

func init() {
	GlobalObject = &GlobalObj{
		Name:    "TiAlloy Client",
		Version: "v1.0.0",

		ServerHost: "127.0.0.1",
		ServerPort: 8888,

		TcpMaxPacketSize:    4048,
		TcpWorkerPoolSize:   2,
		TcpMaxWorkerTaskLen: 2,

		ReconnectInterval: 1,

		LogMode: true, // true：详细，打印log在代码中输出位置；false：简要，不打印文件输出位置，不打印debug和trace（性能高，生产环境使用）
	}

	GlobalObject.Reload()
}
