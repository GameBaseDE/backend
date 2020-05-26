package openapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpRequestParser struct {
	nextHandler httpRequestHandler
}

func newHttpRequestParser() *httpRequestParser {
	return &httpRequestParser{nextHandler: newHttpRequestKubernetesTranslator()}
}

// ListTemplates - Get a list of all available game server images
func (hr *httpRequestParser) ListTemplates(c *gin.Context) {
	//no parameter checks for list
	hr.nextHandler.ListTemplates(c)
}

// GetStatus - Query status of all deployments
func (hr *httpRequestParser) GetStatus(c *gin.Context) {
	//no parameter checks for status
	hr.nextHandler.GetStatus(c)
}

// ConfigureContainer - Configure a game server based on POST body
func (hr *httpRequestParser) ConfigureContainer(c *gin.Context) {
	//test for parameters
	hr.nextHandler.ConfigureContainer(c)
}

// DeployContainer - Deploy a game server based on POST body
func (hr *httpRequestParser) DeployContainer(c *gin.Context) {
	var request GameContainerDeployment
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hr.nextHandler.DeployContainer(c)
}

// StartContainer - Start a game server/container
func (hr *httpRequestParser) StartContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}
	hr.nextHandler.StartContainer(c)
}

// StopContainer - Stop a game server/container
func (hr *httpRequestParser) StopContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}
	hr.nextHandler.StopContainer(c)
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestParser) RestartContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}
	hr.nextHandler.RestartContainer(c)
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestParser) DeleteContainer(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		id = c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
			return
		}
	}
	hr.nextHandler.DeleteContainer(c)
}
