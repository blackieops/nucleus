package nxc

import (
	"fmt"
	"net/http"
	"crypto/md5"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/files"
	"github.com/gin-gonic/gin"
)

type NextcloudRouter struct {
	DBContext      *data.Context
	Config         *config.Config
	StorageBackend files.StorageBackend
}

func (n *NextcloudRouter) Mount(r *gin.RouterGroup) {
	middleware := &auth.AuthMiddleware{Config: n.Config}

	r.GET("/status.php", func(c *gin.Context) {
		payload := &StatusResponse{
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
		session := CreateNextcloudAuthSession(n.DBContext)
		payload := &PollResponse{
			Poll: PollEndpoint{
				Token:       session.PollToken,
				EndpointURL: n.Config.BaseURL + "/nextcloud/index.php/login/v2/poll",
			},
			LoginURL: n.Config.BaseURL + "/nextcloud/index.php/login/v2/grant?token=" + session.LoginToken,
		}

		c.JSON(201, payload)
	})

	r.POST("/index.php/login/v2/poll", func(c *gin.Context) {
		token := c.PostForm("token")
		session, err := FindNextcloudAuthSessionByPollToken(n.DBContext, token)
		if err != nil || session.RawAppPassword == "" {
			c.JSON(404, make([]string, 0))
			return
		}
		payload := &PollSuccessResponse{
			Server:   n.Config.BaseURL + "/nextcloud",
			Username: session.Username,
			Password: session.RawAppPassword,
		}
		DestroyNextcloudAuthSession(n.DBContext, session)
		c.JSON(200, payload)
	})

	r.GET("/index.php/login/v2/grant", middleware.EnsureSession, func(c *gin.Context) {
		token := c.Query("token")
		c.HTML(200, "nextcloud_grant.html", gin.H{"Token": token})
	})

	r.POST("/index.php/login/v2/grant", middleware.EnsureSession, func(c *gin.Context) {
		user := auth.CurrentUser(n.DBContext, c)
		authSession, err := FindNextcloudAuthSessionByLoginToken(n.DBContext, c.Query("token"))
		if err != nil {
			c.JSON(404, gin.H{"error": err})
			return
		}
		CreateNextcloudAppPassword(n.DBContext, authSession, user)
		c.HTML(201, "nextcloud_grant_success.html", gin.H{})
	})

	r.GET("/ocs/v1.php/cloud/capabilities", func(c *gin.Context) {
		c.JSON(200, BuildCapabilitiesResponse())
	})

	r.GET("/remote.php/dav/avatars/:username/:size.png", func(c *gin.Context) {
		user, err := CurrentUser(n.DBContext, c)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		gravatarUrl := fmt.Sprintf(
			"https://www.gravatar.com/avatar/%x?s=%s",
			md5.Sum([]byte(user.EmailAddress)),
			c.Params.ByName("size"),
		)
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
	})

	webdavRouter := &WebdavRouter{DBContext: n.DBContext, Backend: n.StorageBackend}

	r.Handle("PROPFIND", "/remote.php/dav/files/:username/*filePath", webdavRouter.HandlePropfind)
	//r.Handle("PROPPATCH", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("GET", "/remote.php/dav/files/:username/*filePath", webdavRouter.HandleGet)
	r.Handle("PUT", "/remote.php/dav/files/:username/*filePath", webdavRouter.HandlePut)
	r.Handle("MKCOL", "/remote.php/dav/files/:username/*filePath", webdavRouter.HandleMkcol)

	r.Handle("MKCOL", "/remote.php/dav/uploads/:username/*filePath", webdavRouter.HandleChunkMkcol)
	r.Handle("PUT", "/remote.php/dav/uploads/:username/*filePath", webdavRouter.HandleChunkPut)
	r.Handle("MOVE", "/remote.php/dav/uploads/:username/*filePath", webdavRouter.HandleChunkMove)

	// TODO: Copy? Move? Delete? Others??
}
