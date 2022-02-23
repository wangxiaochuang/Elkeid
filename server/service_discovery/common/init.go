package common

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common/ylog"
)

var (
	Sig  = make(chan os.Signal, 1)
	Quit = make(chan bool)
	V    *viper.Viper

	ConfigChangeNotify = make(chan bool, 1)
	RunMode            string
	SrvIp              string
	SrvPort            int

	AuthEnable bool
	AuthKeys   map[string]string
)

const (
	defaultConfigFile = "./conf/conf.yaml"
)

func init() {
	signal.Notify(Sig, os.Interrupt)
	V = viper.New()
	V.SetConfigFile(defaultConfigFile)
	V.WatchConfig()
	V.OnConfigChange(watchCallback)
	err := V.ReadInConfig()
	if err != nil {
		fmt.Printf("read in config error: %s\n", err.Error())
		os.Exit(1)
	}
	configNotify()
	logger := ylog.NewYLog(
		ylog.WithLogFile(V.GetString("log.path")),
		ylog.WithMaxAge(3),
		ylog.WithMaxSize(10),
		ylog.WithMaxBackups(3),
		ylog.WithLevel(V.GetInt("log.loglevel")),
	)
	ylog.InitLogger(logger)

	RunMode = V.GetString("Cluster.Mode")
	SrvIp = V.GetString("Server.Ip")
	SrvPort = V.GetInt("Server.Port")

	AuthEnable = V.GetBool("Auth.Enable")
	AuthKeys = V.GetStringMapString("Auth.Keys")
	return
}

func configNotify() {
	select {
	case ConfigChangeNotify <- true:
		fmt.Printf("notify config changed\n")
	default:
		fmt.Printf("config change notify channel is block\n")
	}
}

func watchCallback(e fsnotify.Event) {
	configNotify()
}
