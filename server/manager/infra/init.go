package infra

import (
	"flag"
	"fmt"
	"os"

	"github.com/wangxiaochuang/Elkeid/server/manager/infra/userconfig"
	"github.com/wangxiaochuang/Elkeid/server/ylog"
)

func init() {
	confPath := flag.String("c", "conf/svr.yml", "ConfigPath")
	flag.Parse()
	ConfPath = *confPath

	InitConfig()
}

func initlog() {
	logger := ylog.NewYLog(
		ylog.WithLogFile(Conf.GetString("log.path")),
		ylog.WithMaxAge(3),
		ylog.WithMaxSize(10),
		ylog.WithMaxBackups(3),
		ylog.WithLevel(Conf.GetInt("log.loglevel")),
	)
	ylog.InitLogger(logger)
}

func InitConfig() {
	var (
		err error
	)
	if Conf, err = userconfig.NewUserConfig(userconfig.WithPath(ConfPath)); err != nil {
		fmt.Println("NEW_CONFIG_ERROR", err.Error())
		os.Exit(-1)
	}

	initlog()
	// initComponents()
	// initDefault()
}
