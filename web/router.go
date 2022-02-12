package web

import (
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/internal/csrf"
	"github.com/gin-gonic/gin"
)

type WebRouter struct {
	DBContext *data.Context
}

func (r *WebRouter) Mount(g *gin.RouterGroup) {
	g.GET("/login", csrf.Generate(), r.handleLoginShow)
	g.POST("/login", csrf.Validate(), r.handleLoginCreate)
}
