package main

import (
	"os/signal"
	"syscall"

	"github.com/wangxiaochuang/Elkeid/server/manager/infra"
	"github.com/wangxiaochuang/Elkeid/server/manager/infra/discovery"
)

func init() {
	signal.Notify(infra.Sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

func main() {
	reg := discovery.NewServerRegistry()
	defer reg.Stop()

	go ServerStart()

	<-infra.Sig
	close(infra.Quit)
}

func ServerStart() {

}
