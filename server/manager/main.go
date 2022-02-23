package main

import (
	"os/signal"
	"syscall"

	"github.com/wangxiaochuang/Elkeid/server/manager/infra"
)

func init() {
	signal.Notify(infra.Sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

func main() {}
