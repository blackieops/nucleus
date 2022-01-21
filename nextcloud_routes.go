package main

import (
	"github.com/gin-gonic/gin"
	"com.blackieops.nucleus/nxc"
)

func mountNextcloudRoutes(r *gin.RouterGroup) {
	r.GET("/status.php", func (c *gin.Context) {
		payload := &nxc.StatusResponse{
			Installed: true,
			Maintenance: false,
			NeedsDatabaseUpgrade: false,
			Version: "22.2.3.0",
			VersionString: "22.2.3",
			Edition: "",
			ProductName: "Nextcloud",
			ExtendedSupport: false,
		}
		c.JSON(200, payload)
	})

	r.POST("/index.php/login/v2", func (c *gin.Context) {
		// TODO: implement login v2 flow
	})

	r.Handle("PROPFIND", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("PROPPATCH", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("GET", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	r.Handle("PUT", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)
	// TODO: Copy? Move? Delete? Others??
}
