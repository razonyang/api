package main

import (
	"os"

	"github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/razonyang/api/internal/app"
	"github.com/razonyang/api/internal/github"
	"github.com/razonyang/api/internal/hugo"
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

	cacheStore := redis_store.NewRedis(redis.NewClient(app.RedisOptions()))
	cacheManager := cache.New[string](cacheStore)

	githubCtrl := github.New(github.NewService(), cacheManager)
	moduleCtrl := hugo.NewModuleController(hugo.NewService(cacheManager))

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})
	v1 := r.Group("/v1")
	v1.GET("/github/dependents/:owner/:repo", githubCtrl.Dependents)
	v1.GET("/hugo/modules/:vendor/:owner/:repo", moduleCtrl.Requirements)
	r.Run()
}
