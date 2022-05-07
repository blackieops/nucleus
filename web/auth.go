package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.b8s.dev/nucleus/auth"
)

func (r *WebRouter) handleLoginShow(c *gin.Context) {
	s := sessions.Default(c)
	c.HTML(200, "auth_login.html", gin.H{"csrfToken": s.Get("CSRFToken")})
}

func (r *WebRouter) handleLoginCreate(c *gin.Context) {
	user, err := auth.FindUserByUsername(r.DBContext, c.PostForm("username"))
	if err != nil {
		// We still run through a "fake" credential validation to prevent
		// leaking credential/user existence through response time
		// variance.
		auth.ValidateCredential(&auth.Credential{}, "burnsometime")
		c.AbortWithStatus(404)
		return
	}
	credentials, err := auth.FindUserCredentials(r.DBContext, user)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	credential, err := auth.FilterFirstCredentialOfType(credentials, auth.CredentialTypePassword)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	err = auth.ValidateCredential(credential, c.PostForm("password"))
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
}

func (r *WebRouter) handleLoginDestroy(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(302, "/")
}
