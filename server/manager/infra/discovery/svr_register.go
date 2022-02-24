package discovery

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/levigross/grequests"
	"github.com/wangxiaochuang/Elkeid/server/manager/infra"
	"github.com/wangxiaochuang/Elkeid/server/manager/infra/biz/midware"
	"github.com/wangxiaochuang/Elkeid/server/ylog"
)

const (
	sdRegisterURL           = "http://%s/registry/register"
	sdEvictURL              = "http://%s/registry/evict"
	defaultRegisterInterval = 30
)

type ServerRegistry struct {
	Name     string `json:"name"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Weight   int    `json:"weight"`
	SDHost   string
	stopChan chan struct{}
}

func NewServerRegistry() *ServerRegistry {
	host := infra.SDAddrs[rand.Int()%len(infra.SDAddrs)]
	return NewRegistry(infra.RegisterName, infra.LocalIP, host, infra.HttpPort)
}

func tmpFuncs(svr *ServerRegistry, host string) (string, error) {
	option := midware.SdAuthRequestOption()
	option.JSON = map[string]interface{}{
		"name":   svr.Name,
		"ip":     svr.Ip,
		"port":   svr.Port,
		"weight": svr.Weight,
	}
	option.RequestTimeout = 2 * time.Second
	url := fmt.Sprintf(sdRegisterURL, host)
	r, err := grequests.Post(url, option)
	if err != nil {
		return "", err
	}
	return r.String(), err
}

func _request(url string, data map[string]interface{}) (string, error) {
	option := midware.SdAuthRequestOption()
	option.JSON = data
	option.RequestTimeout = 2 * time.Second
	r, err := grequests.Post(url, option)
	if err != nil {
		return "", err
	}
	return r.String(), err
}

func NewRegistry(svrName, ip, sdUrl string, port int) *ServerRegistry {
	svr := &ServerRegistry{
		Name:     svrName,
		Ip:       ip,
		Port:     port,
		Weight:   0,
		SDHost:   sdUrl,
		stopChan: make(chan struct{}),
	}

	ylog.Infof("NewRegistry", "new registry: %#v", *svr)
	r, err := tmpFuncs(svr, sdUrl)
	if err != nil {
		ylog.Errorf("NewRegistry", "register failed: %v", err)
		fmt.Printf("register error: %s\n", err.Error())
		return svr
	}
	ylog.Infof("NewRegistry", "register response: %s", r)
	go svr.renewRegistry()
	return svr
}

func (s *ServerRegistry) renewRegistry() {
	t := time.NewTicker(defaultRegisterInterval * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			resp, err := tmpFuncs(s, s.SDHost)
			if err != nil {
				ylog.Errorf("RenewRegistry", "####renew registry failed: %v", err)
				continue
			}
			ylog.Debugf("RenewRegistry", ">>>>renew registry resp: %s", resp)
		case <-s.stopChan:
			return
		}
	}
}

func (s *ServerRegistry) Stop() {
	var (
		err  error
		resp string
	)
	close(s.stopChan)

	url := fmt.Sprintf(sdEvictURL, s.SDHost)
	data := map[string]interface{}{
		"name": s.Name,
		"ip":   s.Ip,
		"port": s.Port,
	}
	resp, err = _request(url, data)
	if err != nil {
		ylog.Errorf("ServerRegistryStop", "####evict server failed: %v", err)
		return
	}

	ylog.Debugf("ServerRegistryStop", ">>>>evict server resp: %s", resp)
}

func (s *ServerRegistry) SetWeight(w int) {
	s.Weight = w
}
