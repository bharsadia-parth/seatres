package api

import (
	"github.com/bharsadia-parth/seatres/internal/store"
	"github.com/gin-gonic/gin"
)

func NewRouter(st *store.Store, adminToken string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), RequestID(), AccessLog())
	r.GET("/healthz", Healthz)
	r.GET("/readyz", Readyz(st))
	r.GET("/shows/:id", GetShow(st))

	authed := r.Group("/", Auth(adminToken))
	authed.POST("/shows", AdminOnly(), CreateShow(st))
	return r
}
