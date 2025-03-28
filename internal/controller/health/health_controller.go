package health

import (
	"net/http"
	"time"

	config "battledak-server/configs"

	"github.com/gin-gonic/gin"
)

type HealthController interface {
	Check(c *gin.Context)
}

type healthController struct{}

func NewHealthController() HealthController {
	return &healthController{}
}

func (h *healthController) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"env":       config.GetEnv("ENV"),
		"version":   "1.0.0",
	})
}
