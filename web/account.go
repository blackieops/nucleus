package web

import (
	"com.blackieops.nucleus/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (r *WebRouter) handleAccountEdit(c *gin.Context) {
	s := sessions.Default(c)
	user := r.Auth.GetCurrentUser(c)
	csrfToken := s.Get("CSRFToken")
	c.HTML(200, "account_edit.html", gin.H{"user": user, "csrfToken": csrfToken})
}

func (r *WebRouter) handleAccountUpdate(c *gin.Context) {
	user := r.Auth.GetCurrentUser(c)
	if newEmail := c.PostForm("emailAddress"); newEmail != "" {
		user.EmailAddress = newEmail
	}
	user, err := auth.UpdateUser(r.DBContext, user)
	if err != nil {
		panic(err)
	}
	if newPassword := c.PostForm("password"); newPassword != "" {
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
		_, err = auth.UpdateCredential(r.DBContext, credential, newPassword)
		if err != nil {
			c.AbortWithStatus(422)
			return
		}
	}
	c.Redirect(302, "/web/")
}
