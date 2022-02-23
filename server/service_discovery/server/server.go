package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common/ylog"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/server/handler"
)

func ServerStart(ip string, port int) {
	r := gin.Default()
	register(r)
	go func() {
		ylog.Infof("[START_SERVER]", "Listening and serving on :%s:%d\n", ip, port)
		fmt.Printf("server run error: %s\n", r.Run(fmt.Sprintf("%s:%d", ip, port)).Error())
	}()

	select {
	case <-common.Sig:
		handler.EI.Stop()
		handler.CI.Stop()
		close(common.Quit)
	}
}
