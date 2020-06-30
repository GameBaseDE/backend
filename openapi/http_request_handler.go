package openapi

import (
	"github.com/gin-gonic/gin"
)

type httpRequestHandler interface {
	kubernetesClient() kubernetesClient
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)
	ListTemplates(c *gin.Context)
	GetStatus(c *gin.Context)
	ConfigureContainer(c *gin.Context)
	DeployContainer(c *gin.Context)
	StartContainer(c *gin.Context)
	StopContainer(c *gin.Context)
	RestartContainer(c *gin.Context)
	DeleteContainer(c *gin.Context)
}
