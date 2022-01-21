package main

import (
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/nxc"
	"github.com/gin-gonic/gin"
)

type NextcloudRouter struct {
	DBContext *data.Context
	Config    *config.Config
}

func (n *NextcloudRouter) Mount(r *gin.RouterGroup) {
	middleware := &MiddlewareRouter{Config: n.Config}

	r.GET("/status.php", func(c *gin.Context) {
		payload := &nxc.StatusResponse{
			Installed:            true,
			Maintenance:          false,
			NeedsDatabaseUpgrade: false,
			Version:              "22.2.3.0",
			VersionString:        "22.2.3",
			Edition:              "",
			ProductName:          "Nextcloud",
			ExtendedSupport:      false,
		}
		c.JSON(200, payload)
	})

	r.POST("/index.php/login/v2", func(c *gin.Context) {
		session := data.CreateNextcloudAuthSession(n.DBContext)
		payload := &nxc.PollResponse{
			Poll: nxc.PollEndpoint{
				Token:       session.PollToken,
				EndpointURL: n.Config.BaseURL + "/nextcloud/index.php/login/v2/poll",
			},
			LoginURL: n.Config.BaseURL + "/nextcloud/index.php/login/v2/grant?token=" + session.LoginToken,
		}

		c.JSON(201, payload)
	})

	r.POST("/index.php/login/v2/poll", func(c *gin.Context) {
		token := c.PostForm("token")
		session, err := data.FindNextcloudAuthSessionByPollToken(n.DBContext, token)
		if err != nil || session.RawAppPassword == "" {
			c.JSON(404, make([]string, 0))
			return
		}
		payload := &nxc.PollSuccessResponse{
			Server:   n.Config.BaseURL,
			Username: session.Username,
			Password: session.RawAppPassword,
		}
		data.DestroyNextcloudAuthSession(n.DBContext, session)
		c.JSON(200, payload)
	})

	r.GET("/index.php/login/v2/grant", middleware.EnsureSession, func(c *gin.Context) {
		token := c.Query("token")
		c.HTML(200, "nextcloud_grant.html", gin.H{"Token": token})
	})

	r.POST("/index.php/login/v2/grant", middleware.EnsureSession, func(c *gin.Context) {
		user := currentUser(n.DBContext, c)
		authSession, err := data.FindNextcloudAuthSessionByLoginToken(n.DBContext, c.Query("token"))
		if err != nil {
			c.JSON(404, gin.H{"error": err})
			return
		}
		data.CreateNextcloudAppPassword(n.DBContext, authSession, user)
		c.HTML(201, "nextcloud_grant_success.html", gin.H{})
	})

	r.Handle("PROPFIND", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("PROPPATCH", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("GET", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("PUT", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	// TODO: Copy? Move? Delete? Others??
}
