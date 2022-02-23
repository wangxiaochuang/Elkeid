package handler

import (
	"fmt"

	"github.com/wangxiaochuang/Elkeid/server/service_discovery/cluster"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/endpoint"
)

var (
	CI cluster.Cluster
	EI *endpoint.Endpoint
)

func init() {
	if common.RunMode == "config" {
		CI = cluster.NewConfigCluster(fmt.Sprintf("%s:%d", common.SrvIp, common.SrvPort))
		EI = endpoint.NewEndpoint(CI)
	}
}
