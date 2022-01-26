package main

import (
	"flag"
	"fmt"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/nxc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file.")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	dbContext := data.Connect(conf.DatabaseURL)
	auth.AutoMigrate(dbContext)
	nxc.AutoMigrate(dbContext)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	sessionStore := cookie.NewStore([]byte(conf.SessionSecret))
	r.Use(sessions.Sessions("nucleussession", sessionStore))

	nextcloudRouter := &nxc.NextcloudRouter{
		DBContext: dbContext,
		Config:    conf,
	}
	nextcloudRouter.Mount(r.Group("/nextcloud"))

	authRouter := &auth.AuthRouter{DBContext: dbContext}
	authRouter.Mount(r.Group("/auth"))

	r.Run(fmt.Sprintf(":%d", conf.Port))
}
