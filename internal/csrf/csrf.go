package csrf

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Generate() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set("CSRFToken", generateToken())
		s.Save()
	}
}

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
