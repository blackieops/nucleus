package main

import (
	"flag"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/config"
	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/nxc"
	"com.blackieops.nucleus/files"
	"github.com/gin-contrib/sessions/cookie"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to configuration file.")
	wantIndex = flag.Bool("index", false, "Index the user files on-disk instead of running the server.")
	wantMigrate = flag.Bool("migrate", false, "Run database migrations instead of running the server.")
	wantSeeds = flag.Bool("seed", false, "Insert test data into the database.")
)

func main() {
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

	if *wantSeeds {
		seedData(dbContext)
		return
	}

	fsBackend := &files.FilesystemBackend{StoragePrefix: conf.DataPath}

	if *wantIndex {
		(&files.Crawler{DBContext: dbContext, Backend: fsBackend}).ReindexAll()
		return
	}

	r := &NucleusRouter{
		Auth:           &auth.AuthMiddleware{DBContext: dbContext, Config: conf},
		Config:         conf,
		DBContext:      dbContext,
		SessionStore:   cookie.NewStore([]byte(conf.SessionSecret)),
		StorageBackend: fsBackend,
	}
	r.Configure()
	r.Listen(conf.Port)
}
