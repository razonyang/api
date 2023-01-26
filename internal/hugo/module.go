package hugo

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		log.Panic("unsupported vendor: %s", vendor)
	}

	cfg, err := ctrl.service.Config(c.Request.Context(), vendor, c.Param("owner"), c.Param("repo"))
	if err != nil {
		log.Panic("failed to fetch config: %s", err)
	}

	c.JSON(200, cfg.Module)
}
