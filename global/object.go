package global

import (
	"encoding/json"
	"github.com/kanyuanzhi/tialloy-client/ticlog"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

var Object *Obj

type Obj struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`

	ServerHost string `json:"server_host,omitempty"`
	ServerPort int    `json:"server_port,omitempty"`

	TcpMaxPacketSize    uint32 `json:"tcp_max_packet_size,omitempty"`
	TcpWorkerPoolSize   uint32 `json:"tcp_worker_pool_size,omitempty"`
	TcpMaxWorkerTaskLen uint32 `json:"tcp_max_worker_task_len,omitempty"`

	ReconnectInterval time.Duration `json:"reconnect_interval,omitempty"`

	LogMode bool `json:"log_mode,omitempty"`
}

func (o *Obj) Reload() {
	data, err := ioutil.ReadFile("conf/tialloy_client.json")
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &Object)
	if err != nil {
		panic(err.Error())
	}

	ticlog.Log = logrus.New()
	ticlog.Log.SetReportCaller(Object.LogMode)
	if Object.LogMode == true {
		ticlog.Log.SetLevel(logrus.TraceLevel)
	} else {
		ticlog.Log.SetLevel(logrus.InfoLevel)
	}
	ticlog.Log.SetFormatter(&ticlog.CustomFormatter{})
}

func init() {
	Object = &Obj{
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

	Object.Reload()
}
