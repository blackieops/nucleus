package main

import (
	"flag"
	"fmt"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/files"
	"com.blackieops.nucleus/nxc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file.")
	wantIndex := flag.Bool("index", false, "Index the user files on-disk instead of running the server.")
	wantMigrate := flag.Bool("migrate", false, "Run database migrations instead of running the server.")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	dbContext := data.Connect(conf.DatabaseURL)

	if *wantMigrate {
		auth.AutoMigrate(dbContext)
		nxc.AutoMigrate(dbContext)
		files.AutoMigrate(dbContext)
		return
	}

	fsBackend := &files.FilesystemBackend{StoragePrefix: conf.DataPath}
	if *wantIndex {
		(&files.Crawler{DBContext: dbContext, Backend: fsBackend}).ReindexAll()
		return
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	sessionStore := cookie.NewStore([]byte(conf.SessionSecret))
	r.Use(sessions.Sessions("nucleussession", sessionStore))

	nextcloudRouter := &nxc.NextcloudRouter{
		DBContext:      dbContext,
		Config:         conf,
		StorageBackend: fsBackend,
	}
	nextcloudRouter.Mount(r.Group("/nextcloud"))

	authRouter := &auth.AuthRouter{DBContext: dbContext}
	authRouter.Mount(r.Group("/auth"))

	r.Run(fmt.Sprintf(":%d", conf.Port))
}
