package web

import (
	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
	"go.b8s.dev/nucleus/internal/csrf"
	"github.com/gin-gonic/gin"
)

type WebRouter struct {
	DBContext *data.Context
	Auth      *auth.AuthMiddleware
}

func (r *WebRouter) Mount(g *gin.RouterGroup) {
	g.GET("/login", csrf.Generate(), r.handleLoginShow)
	g.POST("/login", csrf.Validate(), r.handleLoginCreate)
	g.POST("/logout", r.Auth.EnsureSession, csrf.Validate(), r.handleLoginDestroy)

	g.GET("/", r.Auth.EnsureSession, r.handleDashboardShow)

	g.GET("/me/avatar", r.Auth.EnsureSession, r.handleUserAvatarShow)
	g.GET("/me", r.Auth.EnsureSession, csrf.Generate(), r.handleAccountEdit)
	g.POST("/me", r.Auth.EnsureSession, csrf.Validate(), r.handleAccountUpdate)
	g.POST("/me/revokeAppPassword", r.Auth.EnsureSession, csrf.Validate(), r.handleAccountRevokeAppPassword)

	g.GET("/users", r.Auth.EnsureSession, csrf.Generate(), r.handleUsersIndex)
	g.GET("/users/new", r.Auth.EnsureSession, csrf.Generate(), r.handleUsersNew)
	g.POST("/users/new", r.Auth.EnsureSession, csrf.Validate(), r.handleUsersCreate)
	g.POST("/users/:userId/destroy", r.Auth.EnsureSession, csrf.Validate(), r.handleUsersDestroy)
}
