package web

import (
	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/internal/csrf"
	"github.com/gin-gonic/gin"
)

type WebRouter struct {
	DBContext *data.Context
	Auth      *auth.AuthMiddleware
}

func (r *WebRouter) Mount(g *gin.RouterGroup) {
	g.GET("/login", csrf.Generate(), r.handleLoginShow)
	g.POST("/login", csrf.Validate(), r.handleLoginCreate)

	g.GET("/", r.Auth.EnsureSession, csrf.Generate(), r.handleAccountEdit)
	g.POST("/", r.Auth.EnsureSession, csrf.Validate(), r.handleAccountUpdate)
}
