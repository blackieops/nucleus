package nxc

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
)

type Middleware struct {
	DBContext *data.Context
}

// Finds the user making the given request, based on the Nextcloud App Password
// being used in the Authorization header. This is a Nextcloud-specific auth
// strategy and should only be used on endpoints that need Nextcloud client
// compatibility.
func (m *Middleware) EnsureAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderRaw := c.GetHeader("Authorization")
		if authHeaderRaw == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		authHeaderValue := authHeaderRaw[len("Basic "):]
		authHeaderValueBytes, err := base64.URLEncoding.DecodeString(authHeaderValue)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		authHeaderPassword := strings.Split(string(authHeaderValueBytes), ":")[1]
		appPassword, err := FindNextcloudAppPasswordByPassword(m.DBContext, authHeaderPassword)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("CurrentUser", &appPassword.User)
	}
}

func (r *Middleware) GetCurrentUser(c *gin.Context) *auth.User {
	user, exist := c.Get("CurrentUser")
	if !exist {
		panic("You need to add EnsureAuthorization as middleware before calling GetCurrentUser.")
	}
	return user.(*auth.User)
}
