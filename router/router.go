package router

import (
	"github.com/gin-gonic/gin"

	"github.com/seeu359/go-project-278/handlers"
)

func New() *gin.Engine {
	r := gin.Default()
	r.GET("ping", handlers.Ping)

	return r
}
