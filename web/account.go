package web

import (
	"fmt"
	"strconv"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/nxc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (r *WebRouter) handleAccountEdit(c *gin.Context) {
	s := sessions.Default(c)
	user := r.Auth.GetCurrentUser(c)
	nxcPasswords, _ := nxc.ListNextcloudAppPasswordsForUser(r.DBContext, user)
	csrfToken := s.Get("CSRFToken")
	c.HTML(200, "account_edit.html", gin.H{"user": user, "csrfToken": csrfToken, "appPasswords": nxcPasswords})
}

func (r *WebRouter) handleAccountUpdate(c *gin.Context) {
	user := r.Auth.GetCurrentUser(c)
	if newEmail := c.PostForm("emailAddress"); newEmail != "" {
		user.EmailAddress = newEmail
	}
	if newName := c.PostForm("name"); newName != "" {
		user.Name = newName
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

func (r *WebRouter) handleAccountRevokeAppPassword(c *gin.Context) {
	user := r.Auth.GetCurrentUser(c)
	s := sessions.Default(c)
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		s.AddFlash(fmt.Sprintf("Could not revoke Nextcloud App Password: %v", err))
		c.Redirect(302, "/web/me")
		return
	}
	password, err := nxc.FindNextcloudAppPassword(r.DBContext, user, uint(id))
	if err != nil {
		s.AddFlash(fmt.Sprintf("Could not revoke Nextcloud App Password: %v", err))
		c.Redirect(302, "/web/me")
		return
	}
	err = nxc.DeleteNextcloudAppPassword(r.DBContext, password)
	if err != nil {
		s.AddFlash(fmt.Sprintf("Could not revoke Nextcloud App Password: %v", err))
		c.Redirect(302, "/web/me")
		return
	}
	c.Redirect(302, "/web/me")
}
