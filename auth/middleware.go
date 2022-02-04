package auth

import (
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Config *config.Config
}

// Middleware to check if there is a currently logged-in user in the session.
func (r *AuthMiddleware) EnsureSession(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("CurrentUserID") == nil {
		session.Set("ReturnTo", c.Request.URL.Path+"?"+c.Request.URL.RawQuery)
		session.Save()
		c.Redirect(302, r.Config.BaseURL+"/auth/login")
		c.Abort()
		return
	}
}

func CurrentUser(c *data.Context, g *gin.Context) (*User, error) {
	session := sessions.Default(g)
	userID := session.Get("CurrentUserID").(uint)
	return FindUser(c, userID)
}
