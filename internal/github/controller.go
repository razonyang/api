package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	service *Service
	cache   cache.CacheInterface[string]
}

func New(service *Service, cache cache.CacheInterface[string]) *Controller {
	return &Controller{
		service: service,
		cache:   cache,
	}
}

func (controller *Controller) Dependents(c *gin.Context) {
	cacheKey := fmt.Sprintf("github-dependents -%s-%s", c.Param("owner"), c.Param("repo"))
	ctx := c.Request.Context()
	dependents := &Dependents{}
	cachedVal, err := controller.cache.Get(ctx, cacheKey)
	if err == nil {
		err = json.Unmarshal([]byte(cachedVal), dependents)
	}
	if err != nil {
		log.Debugf("failed to fetch dependents from cache: %s", err)

		dependents, err = controller.service.DependentsByOwnerAndRepo(c.Param("owner"), c.Param("repo"))
		if err != nil {
			log.Panicf("failed to fetch dependents: %s", err)
		}

		data, err := json.Marshal(dependents)
		if err != nil {
			log.Infof("failed to marshal dependents: %s", err)
		} else {
			if err = controller.cache.Set(ctx, cacheKey, string(data), store.WithExpiration(6*time.Hour)); err != nil {
				log.Infof("failed to cache: %s", err)
			}
		}
	}
	fmt.Println(dependents)
	c.JSON(http.StatusOK, dependents)
}
