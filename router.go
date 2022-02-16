package main

import (
	"fmt"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/files"
	"com.blackieops.nucleus/nxc"
	"com.blackieops.nucleus/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type NucleusRouter struct {
	DBContext      *data.Context
	Config         *config.Config
	SessionStore   sessions.Store
	Auth           *auth.AuthMiddleware
	StorageBackend files.StorageBackend
	router         *gin.Engine
}

func (nr *NucleusRouter) Configure() {
	nr.router = gin.Default()
	nr.router.LoadHTMLGlob("templates/*")
	nr.router.Use(static.Serve("/static", static.LocalFile("static", false)))
	nr.router.Use(sessions.Sessions("nucleussession", nr.SessionStore))

	nextcloudRouter := &nxc.NextcloudRouter{
		DBContext:      nr.DBContext,
		Config:         nr.Config,
		StorageBackend: nr.StorageBackend,
		Auth:           nr.Auth,
	}
	nextcloudRouter.Mount(nr.router.Group("/nextcloud"))

	webRouter := &web.WebRouter{
		DBContext: nr.DBContext,
		Auth:      nr.Auth,
	}
	webRouter.Mount(nr.router.Group("/web"))

	nr.router.GET("/", func(c *gin.Context) {
		// If you hit the root path, you probably wanted the web app.
		c.Redirect(302, "/web/")
	})
}

func (nr *NucleusRouter) Listen(port int) {
	nr.router.Run(fmt.Sprintf(":%d", port))
}
