package http

import (
	"github.com/gin-gonic/gin"
)

type httpRequestHandler interface {
	ListTemplates(c *gin.Context)
	GetStatus(c *gin.Context)
	ConfigureContainer(c *gin.Context)
	DeployContainer(c *gin.Context)
	StartContainer(c *gin.Context)
	StopContainer(c *gin.Context)
	RestartContainer(c *gin.Context)
	DeleteContainer(c *gin.Context)
}
