package auth

import (
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/internal/csrf"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	DBContext *data.Context
}

func (a *AuthRouter) Mount(r *gin.RouterGroup) {
	r.GET("/login", csrf.Generate(), func(c *gin.Context) {
		s := sessions.Default(c)
		c.HTML(200, "auth_login.html", gin.H{"csrfToken": s.Get("CSRFToken")})
	})

	r.POST("/login", csrf.Validate(), func(c *gin.Context) {
		user, err := FindUserByUsername(a.DBContext, c.PostForm("username"))
		if err != nil {
			// We still run through a "fake" credential validation to prevent
			// leaking credential/user existence through response time
			// variance.
			ValidateCredential(&Credential{}, "burnsometime")
			c.AbortWithStatus(404)
			return
		}
		credentials, err := FindUserCredentials(a.DBContext, user)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}
		credential, err := FilterFirstCredentialOfType(credentials, CredentialTypePassword)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		err = ValidateCredential(credential, c.PostForm("password"))
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		session := sessions.Default(c)
		session.Set("CurrentUserID", user.ID)
		session.Save()

		if rt := session.Get("ReturnTo"); rt != nil {
			c.Redirect(302, rt.(string))
			c.Abort()
			return
		}

		c.Redirect(302, "/")
	})
}
