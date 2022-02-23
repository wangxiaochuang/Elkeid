package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/common"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/server/handler"
	"github.com/wangxiaochuang/Elkeid/server/service_discovery/server/midware"
)

func register(r *gin.Engine) {
	authorized := r.Group("/")
	if common.AuthEnable {
		authorized.Use(midware.AKSKAuth())
	}
	{
		authorized.POST("/registry/register", handler.Register)
		authorized.POST("/registry/evict", handler.Evict)
		authorized.POST("/registry/sync", handler.Sync)
	}

	r.GET("/endpoint/ping", handler.Ping)
	r.GET("/endpoint/stat", handler.EndpointStat)

	r.GET("/registry/summary", handler.RegistrySummary)
	r.GET("/registry/detail", handler.RegistryDetail)
	r.GET("/registry/list", handler.RegistryList)
}
