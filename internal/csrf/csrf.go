package csrf

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Generate is a Gin middleware function that generates a CSRF token and adds it
// to the session.
func Generate() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set("CSRFToken", generateToken())
		s.Save()
	}
}

// Validate is a Gin middleware function that checks the `_csrf` post-form value
// with the CSRF token in the session. Aborts the request if they mismatch.
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		token := s.Get("CSRFToken")
		if token == nil {
			c.AbortWithStatus(400)
			return
		}
		paramsToken := c.PostForm("_csrf")
		if token.(string) != paramsToken {
			c.AbortWithStatus(401)
		}
	}
}

func generateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
