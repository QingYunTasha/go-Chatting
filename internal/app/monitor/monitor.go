package monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"liveness": "OK",
	})
}

func ReadinessCheck(c *gin.Context) {
	// check func

	c.JSON(http.StatusOK, gin.H{
		"readiness": "OK",
	})
}

func MetricsExporter(c *gin.Context) {}
