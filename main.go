package main

import (
	"github.com/gin-gonic/gin"
	"com.blackieops.nucleus/webdav"
)

func main() {
	r := gin.Default()

	nextcloudRoutes := r.Group("/nextcloud")
	mountNextcloudRoutes(nextcloudRoutes)

	r.Run(":8989")
}

func forwardToWebdav(c *gin.Context) {
	webdav.HandleRequest(c.Writer, c.Request)
}
