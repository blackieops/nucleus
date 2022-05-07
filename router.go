package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/config"
	"go.b8s.dev/nucleus/data"
	"go.b8s.dev/nucleus/files"
	"go.b8s.dev/nucleus/nxc"
	"go.b8s.dev/nucleus/web"
)

type NucleusRouter struct {
	DBContext      *data.Context
	Config         *config.Config
	SessionStore   sessions.Store
	Auth           *auth.AuthMiddleware
	StorageBackend files.StorageBackend
	router         *gin.Engine
}

//go:embed templates/* static/*
var assetFS embed.FS

func (nr *NucleusRouter) Configure() {
	nr.router = gin.Default()

	// Use binary-embedded templates and static assets.
	tmpls := template.Must(template.New("").ParseFS(assetFS, "templates/*"))
	staticlessAssetFS, err := fs.Sub(assetFS, "static")
	if err != nil {
		panic(err)
	}
	nr.router.SetHTMLTemplate(tmpls)
	nr.router.StaticFS("/static", http.FS(staticlessAssetFS))

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

	nr.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"alive": true})
	})

	nr.router.GET("/", func(c *gin.Context) {
		// If you hit the root path, you probably wanted the web app.
		c.Redirect(302, "/web/")
	})
}

func (nr *NucleusRouter) Listen(port int) {
	nr.router.Run(fmt.Sprintf(":%d", port))
}
