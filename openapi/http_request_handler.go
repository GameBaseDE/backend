package openapi

import (
	"github.com/gin-gonic/gin"
)

type httpRequestHandler interface {
	ListImages(c *gin.Context)
	GetStatus(c *gin.Context)
	ConfigureContainer(c *gin.Context)
	DeployContainer(c *gin.Context)
	StartContainer(c *gin.Context)
	StopContainer(c *gin.Context)
	RestartContainer(c *gin.Context)
	DeleteContainer(c *gin.Context)
}

/**
type httpRequestDispatcher struct {
	nextHandler httpRequestHandler
}

func newHttpRequestDispatcher() *httpRequestDispatcher {
	return &httpRequestDispatcher{nextHandler: newHttpRequestAuthenticator()}
}

// ListImages - Get a list of all available game server images
func (hr *httpRequestDispatcher) ListImages(c *gin.Context) {
}

// GetStatus - Query status of all deployments
func (hr *httpRequestDispatcher) GetStatus(c *gin.Context) {
}

// ConfigureContainer - Configure a game server based on POST body
func (hr *httpRequestDispatcher) ConfigureContainer(c *gin.Context) {
}

// DeployContainer - Deploy a game server based on POST body
func (hr *httpRequestDispatcher) DeployContainer(c *gin.Context) {
}

// StartContainer - Start a game server/container
func (hr *httpRequestDispatcher) StartContainer(c *gin.Context) {
}

// StopContainer - Stop a game server/container
func (hr *httpRequestDispatcher) StopContainer(c *gin.Context) {
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestDispatcher) RestartContainer(c *gin.Context) {
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestDispatcher) DeleteContainer(c *gin.Context) {
}
**/
