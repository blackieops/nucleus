package auth

import (
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	DBContext *data.Context
	Config    *config.Config
}

// Middleware to check if there is a currently logged-in user in the session.
func (r *AuthMiddleware) EnsureSession(c *gin.Context) {
	s := sessions.Default(c)
	if s.Get("CurrentUserID") == nil {
		r.forceLogin(c, s)
		return
	}
	user, err := FindUser(r.DBContext, s.Get("CurrentUserID").(uint))
	if err != nil {
		r.forceLogin(c, s)
	}
	c.Set("CurrentUser", user)
}

func (r *AuthMiddleware) GetCurrentUser(c *gin.Context) *User {
	user, exist := c.Get("CurrentUser")
	if !exist {
		panic("You need to add EnsureSession as middleware before calling GetCurrentUser.")
	}
	return user.(*User)
}

func (r *AuthMiddleware) forceLogin(c *gin.Context, s sessions.Session) {
	s.Set("ReturnTo", c.Request.URL.Path+"?"+c.Request.URL.RawQuery)
	s.Save()
	c.Redirect(302, r.Config.BaseURL+"/web/login")
	c.Abort()
}
