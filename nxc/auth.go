package nxc

import (
	"encoding/base64"
	"strings"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
	"github.com/gin-gonic/gin"
)

type PollResponse struct {
	Poll     PollEndpoint `json:"poll"`
	LoginURL string       `json:"login"`
}

type PollEndpoint struct {
	Token       string `json:"token"`
	EndpointURL string `json:"endpoint"`
}

type PollSuccessResponse struct {
	Server   string `json:"server"`
	Username string `json:"loginName"`
	Password string `json:"appPassword"`
}

func CurrentUser(c *data.Context, g *gin.Context) (*auth.User, error) {
	authHeaderValue := string(g.GetHeader("Authorization"))[len("Basic "):]
	authHeaderValueBytes, err := base64.URLEncoding.DecodeString(authHeaderValue)
	if err != nil {
		panic(err)
	}
	authHeaderPassword := strings.Split(string(authHeaderValueBytes), ":")[1]
	appPassword, err := FindNextcloudAppPasswordByPassword(c, authHeaderPassword)
	if err != nil {
		return &auth.User{}, err
	} else {
		return &appPassword.User, nil
	}
}
