package monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	Usecase *domain.UsecaseRepo
}

func NewMonitorHandler(server *gin.Engine, usecase *domain.UsecaseRepo) {
	mh := MonitorHandler{
		Usecase: usecase,
	}

	server.GET("/liveness", mh.LivenessCheck)
	server.GET("/readiness", mh.ReadinessCheck)
	server.GET("/metrics", mh.MetricsExporter)
}

func (mh *MonitorHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"liveness": "OK",
	})
}

func (mh *MonitorHandler) ReadinessCheck(c *gin.Context) {
	// check func

	c.JSON(http.StatusOK, gin.H{
		"readiness": "OK",
	})
}

func (mh *MonitorHandler) MetricsExporter(c *gin.Context) {}
