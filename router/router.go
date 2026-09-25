package router

import (
	"github.com/gin-gonic/gin"

	"github.com/seeu359/go-project-278/handlers"
	"github.com/seeu359/go-project-278/links"
)

func New(q *links.Queries) *gin.Engine {
	r := gin.Default()
	h := handlers.New(q)

	r.GET("ping", handlers.Ping)
	r.POST("/api/links", h.CreateLink)
	r.GET("/api/links", h.GetLinks)
	r.GET("/api/links/:id", h.GetLinkById)
	r.PUT("/api/links/:id", h.UpdateLink)
	r.DELETE("/api/links/:id", h.DeleteLink)

	r.GET("/r/:code", h.RedirectTo)
	r.GET("/api/link_visits", h.GetVisits)
	return r
}
