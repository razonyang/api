package github

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	service *Service
}

func New(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller *Controller) Dependents(c *gin.Context) {
	ctx := c.Request.Context()
	d, err := controller.service.Dependents(ctx, c.Param("owner"), c.Param("repo"))
	if err != nil {
		log.Panicf("failed to fetch dependents: %s", err)
	}
	c.JSON(http.StatusOK, d)
}
