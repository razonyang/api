package github

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/gin-gonic/gin"
	"github.com/razonyang/apis/internal/github/dependents"
	"github.com/razonyang/apis/internal/helper"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	cache cache.CacheInterface[string]
}

func New(cache cache.CacheInterface[string]) *Controller {
	return &Controller{
		cache: cache,
	}
}

func (controller *Controller) UsedBy(c *gin.Context) {
	cacheKey := fmt.Sprintf("hugo-used-by-%s-%s", c.Param("owner"), c.Param("repo"))
	ctx := c.Request.Context()
	repositories, err := controller.cache.Get(ctx, cacheKey)
	if err != nil {
		log.Debugf("failed to fetch dependents from cache: %s", err)

		dependents, err := dependents.DependentsByOwnerAndRepo(c.Param("owner"), c.Param("repo"))
		if err != nil {
			log.Panicf("failed to fetch dependents: %s", err)
		}

		repositories = helper.FormatInt(dependents.Repositories)
		if err = controller.cache.Set(ctx, cacheKey, repositories, store.WithExpiration(time.Hour)); err != nil {
			log.Infof("failed to cache: %s", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"schemaVersion": 1,
		"namedLogo":     "github",
		"label":         "used by",
		"message":       repositories,
		"color":         "success",
	})
}
