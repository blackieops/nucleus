package main

import (
	"flag"
	"fmt"

	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/webdav"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file.")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	_, err = gorm.Open(postgres.Open(conf.DatabaseURL), &gorm.Config{})
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
