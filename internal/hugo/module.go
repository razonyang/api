package hugo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ModuleController struct {
	service *Service
}

func NewModuleController(service *Service) *ModuleController {
	return &ModuleController{
		service: service,
	}
}

func (ctrl *ModuleController) Requirements(c *gin.Context) {
	vendor := c.Param("vendor")
	if vendor != "github.com" {
		c.Error(fmt.Errorf("unsupported vendor: %s", vendor))
		return
	}

	cfg, err := ctrl.service.Config(c.Request.Context(), vendor, c.Param("owner"), c.Param("repo"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, cfg.Module)
}
