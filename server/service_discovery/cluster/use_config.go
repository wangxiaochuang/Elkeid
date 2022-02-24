package cluster

import (
	"time"

	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common/safemap"
	"github.com/wangxiaochuang/Elkeid/server/ylog"
)

type ConfigCluster struct {
	BaseCluster
}

func NewConfigCluster(host string) Cluster {
	cc := &ConfigCluster{BaseCluster{
		Mode:    configMode,
		Host:    host,
		Members: safemap.NewSafeMap(defaultClusterName),
	}}

	go cc.refresh()
	go cc.ping()

	return cc
}

func (cc *ConfigCluster) refresh() {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		select {
		case changed := <-common.ConfigChangeNotify:
			if changed {
				members := common.V.GetStringSlice("Cluster.Members")
				cc.Members.Del(defaultClusterName)
				for _, host := range members {
					cc.Members.HSet(defaultClusterName, host, "ok")
				}
			}
		case <-t.C:
		case <-cc.Done:
			ylog.Debugf("refresh", "cluster refesh strop\n")
			return
		}
	}
}
