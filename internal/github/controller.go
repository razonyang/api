package github

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func New(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (ctrl *Controller) Dependents(c *gin.Context) {
	ctx := c.Request.Context()
	d, err := ctrl.service.Dependents(ctx, c.Param("owner"), c.Param("repo"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func (ctrl *Controller) Tag(c *gin.Context) {
	ctx := c.Request.Context()
	d, err := ctrl.service.Tag(ctx, c.Param("owner"), c.Param("repo"), c.Query("prefix"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, d)
}
