package main

import (
	"flag"
	"fmt"

	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/webdav"
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
	data.AutoMigrate(dbContext)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	sessionStore := cookie.NewStore([]byte(conf.SessionSecret))
	r.Use(sessions.Sessions("nucleussession", sessionStore))

	nextcloudRouter := &NextcloudRouter{
		DBContext: dbContext,
		Config:    conf,
	}
	nextcloudRouter.Mount(r.Group("/nextcloud"))

	authRouter := &AuthRouter{DBContext: dbContext}
	authRouter.Mount(r.Group("/auth"))

	r.Run(fmt.Sprintf(":%d", conf.Port))
}

func currentUser(c *data.Context, g *gin.Context) *data.User {
	session := sessions.Default(g)
	userID := session.Get("CurrentUserID").(uint)
	return data.FindUser(c, int(userID))
}

func forwardToWebdav(c *gin.Context) {
	webdav.HandleRequest(c.Writer, c.Request)
}
