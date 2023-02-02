package index

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":     os.Getenv("API_NAME"),
		"endpoint": os.Getenv("API_ENDPOINT"),
		"versions": []string{"v1"},
		"source":   "https://github.com/razonyang/api",
	})
}
