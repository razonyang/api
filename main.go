package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/razonyang/api/internal/app"
	"github.com/razonyang/api/internal/github"
	"github.com/razonyang/api/internal/hugo"
	"github.com/razonyang/api/internal/index"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	err := godotenv.Load(".env", ".env.common")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("failed to load .env and .env.common file")
	}

	cacheService := app.NewCacheService()

	githubCtrl := github.New(github.NewService(cacheService))
	moduleCtrl := hugo.NewModuleController(hugo.NewService(cacheService))

	r := gin.Default()
	r.GET("/", index.Index)
	v1 := r.Group("/v1")
	v1.GET("/github/dependents/:owner/:repo", githubCtrl.Dependents)
	v1.GET("/hugo/modules/:vendor/:owner/:repo", moduleCtrl.Requirements)
	r.Run()
}
