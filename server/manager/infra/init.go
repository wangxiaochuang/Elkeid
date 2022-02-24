package infra

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wangxiaochuang/Elkeid/server/manager/infra/mongodb"
	"github.com/wangxiaochuang/Elkeid/server/manager/infra/redis"
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

func initComponents() {
	var err error

	if Grds, err = redis.NewRedisClient(Conf.GetStringSlice("redis.addrs"), Conf.GetString("redis.mastername"), Conf.GetString("redis.passwd")); err != nil {
		fmt.Println("NEW_REDIS_ERROR", err.Error())
		os.Exit(-1)
	}

	err = Grds.Set(context.Background(), "elkeid_manager_test", "test", time.Second).Err()
	if err != nil {
		fmt.Println("REDIS_ERROR", err.Error())
		os.Exit(-1)
	}

	MongoDatabase = Conf.GetString("mongo.dbname")
	if MongoClient, err = mongodb.NewMongoClient(Conf.GetString("mongo.uri")); err != nil {
		fmt.Println("NEW_MONGO_ERROR", err.Error())
		os.Exit(-1)
	}
}

func initDefault() {
	var err error

	LocalIP, err = GetOutboundIP()
	if err != nil {
		ylog.Fatalf("init", "GET_LOCALIP_ERROR: %s Error: %v", LocalIP, err)
	}

	HttpPort = Conf.GetInt("http.port")
	ApiAuth = Conf.GetBool("http.apiauth.enable")
	Secret = Conf.GetString("http.apiauth.secret")
	InnerAuth = Conf.GetStringMapString("http.innerauth")

	SvrName = Conf.GetString("server.name")
	SvrAK = strings.ToLower(Conf.GetString("server.credentials.ak"))
	SvrSK = Conf.GetString("server.credentials.sk")

	SDAddrs = Conf.GetStringSlice("sd.addrs")
	RegisterName = Conf.GetString("sd.name")
	SdAK = strings.ToLower(Conf.GetString("sd.credentials.ak"))
	SdSK = Conf.GetString("sd.credentials.sk")
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
	initComponents()
	initDefault()
}
