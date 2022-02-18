package web

import (
	"net/http"
	"github.com/gin-gonic/gin"
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
