// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/AntonBraer/urlShorter/internal/usecase"
	"github.com/AntonBraer/urlShorter/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Link) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newLinkRoutes(h, t, l)
	}
}
