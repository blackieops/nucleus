package main

import (
	"fmt"
	"flag"

	"github.com/gin-gonic/gin"
	"com.blackieops.nucleus/webdav"
	"com.blackieops.nucleus/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file.")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	nextcloudRoutes := r.Group("/nextcloud")
	mountNextcloudRoutes(nextcloudRoutes)

	r.Run(fmt.Sprintf(":%d", conf.Port))
}

func forwardToWebdav(c *gin.Context) {
	webdav.HandleRequest(c.Writer, c.Request)
}
