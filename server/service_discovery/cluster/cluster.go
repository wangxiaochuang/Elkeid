package cluster

import (
	"fmt"
	"time"

	"github.com/levigross/grequests"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common/safemap"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common/ylog"
)

const (
	defaultClusterName  = "FindYou"
	defaultPingInterval = 5
	pingUrl             = "http://%s/endpoint/ping"
	defaultPingTimeout  = 1
)

const (
	configMode = iota
)

type Cluster interface {
	refresh()
	ping()
	Stop()
	GetHost() string
	GetHosts() []string
	GetOtherHosts() []string
}

type BaseCluster struct {
	Mode    int
	Host    string
	Members *safemap.SafeMap
	Done    chan bool
}

func (bc *BaseCluster) refresh() {}

func (bc *BaseCluster) ping() {
	t := time.NewTicker(defaultPingInterval * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			hosts := bc.Members.HKeys(defaultClusterName)
			for _, host := range hosts {
				if host == bc.Host {
					continue
				}
				url := fmt.Sprintf(pingUrl, host)
				_, err := grequests.Get(url, &grequests.RequestOptions{
					RequestTimeout: defaultPingTimeout * time.Second,
				})
				if err != nil {
					ylog.Errorf("ping", "ping %s error: %s", host, err.Error())
				}
			}
		case <-bc.Done:
			ylog.Debugf("ping", "cluster ping stop")
			return
		}
	}
}

func (bc *BaseCluster) Stop() {
	close(bc.Done)
}

func (bc *BaseCluster) GetHost() string {
	return bc.Host
}

func (bc *BaseCluster) GetHosts() []string {
	hosts := make([]string, 0)
	for _, host := range bc.Members.HKeys(defaultClusterName) {
		hosts = append(hosts, host)
	}
	return hosts
}

func (bc *BaseCluster) GetOtherHosts() []string {
	hosts := make([]string, 0)
	for _, host := range bc.Members.HKeys(defaultClusterName) {
		if host != bc.Host {
			hosts = append(hosts, host)
		}
	}
	return hosts
}
