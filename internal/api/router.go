package api

import "github.com/gin-gonic/gin"

func NewRouter(p Pinger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), RequestID(), AccessLog())
	r.GET("/healthz", Healthz)
	r.GET("/readyz", Readyz(p))
	return r
}
