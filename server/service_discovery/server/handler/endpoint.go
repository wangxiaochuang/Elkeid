package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	r := make(map[string]interface{})
	r["host"] = CI.GetHost()
	r["member"] = CI.GetHosts()
	c.JSON(http.StatusOK, gin.H{"msg": "ok", "data": r})
}

func EndpointStat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok", "data": CI.GetHosts()})
}
