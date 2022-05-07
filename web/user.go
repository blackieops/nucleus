package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.b8s.dev/nucleus/auth"
)

func (r *WebRouter) handleUserAvatarShow(c *gin.Context) {
	user := r.Auth.GetCurrentUser(c)
	gravatarUrl := user.AvatarURL(128)
	response, err := http.Get(gravatarUrl)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="avatar.png"`,
	}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func (r *WebRouter) handleUsersIndex(c *gin.Context) {
	s := sessions.Default(c)
	user := r.Auth.GetCurrentUser(c)
	users := auth.FindAllUsers(r.DBContext)
	csrfToken := s.Get("CSRFToken")
	flashes := s.Flashes()
	s.Save()
	c.HTML(200, "users_index.html", gin.H{"user": user, "users": users, "csrfToken": csrfToken, "flashes": flashes})
}

func (r *WebRouter) handleUsersNew(c *gin.Context) {
	s := sessions.Default(c)
	user := r.Auth.GetCurrentUser(c)
	csrfToken := s.Get("CSRFToken")
	flashes := s.Flashes()
	s.Save()
	c.HTML(200, "users_new.html", gin.H{"user": user, "csrfToken": csrfToken, "flashes": flashes})
}

func (r *WebRouter) handleUsersCreate(c *gin.Context) {
	user := &auth.User{
		EmailAddress: c.PostForm("emailAddress"),
		Username:     c.PostForm("username"),
		Name:         c.PostForm("name"),
	}
	user, err := auth.CreateUser(r.DBContext, user)
	s := sessions.Default(c)
	if err != nil {
		s.AddFlash(fmt.Sprintf("Failed to create user: %v", err))
		s.Save()
		c.Redirect(302, "/web/users/new")
	}
	credential := &auth.Credential{Data: c.PostForm("password")}
	credential, err = auth.CreateCredential(r.DBContext, user, credential)
	if err != nil {
		s.AddFlash(fmt.Sprintf("Failed to create password: %v", err))
		s.Save()
		c.Redirect(302, "/web/users/new")
	}
	c.Redirect(302, "/web/users")
}

func (r *WebRouter) handleUsersDestroy(c *gin.Context) {
	s := sessions.Default(c)
	userId, err := strconv.Atoi(c.Params.ByName("userId"))
	if err != nil {
		s.AddFlash(fmt.Sprintf("Failed to parse user ID: %v", err))
		s.Save()
	}
	user, err := auth.FindUser(r.DBContext, uint(userId))
	if err != nil {
		s.AddFlash(fmt.Sprintf("Failed to find user: %v", err))
		s.Save()
	}
	err = auth.DeleteUser(r.DBContext, user)
	if err != nil {
		s.AddFlash(fmt.Sprintf("Failed to destroy user: %v", err))
		s.Save()
	}
	c.Redirect(302, "/web/users")
}
