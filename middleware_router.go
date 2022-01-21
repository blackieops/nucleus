package main

import (
	"com.blackieops.nucleus/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MiddlewareRouter struct {
	Config *config.Config
}

// Middleware to check if there is a currently logged-in user in the session.
func (r *MiddlewareRouter) EnsureSession(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("CurrentUserID") == nil {
		session.Set("ReturnTo", c.Request.URL.Path+"?"+c.Request.URL.RawQuery)
		session.Save()
		c.Redirect(302, r.Config.BaseURL+"/auth/login")
		c.Abort()
		return
	}
}
