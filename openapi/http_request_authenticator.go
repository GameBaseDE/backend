package openapi

import "github.com/gin-gonic/gin"

type httpRequestAuthenticator struct {
	nextHandler httpRequestHandler
}

func newHttpRequestAuthenticator() *httpRequestAuthenticator {
	return &httpRequestAuthenticator{nextHandler: newHttpRequestParser()}
}

// ListImages - Get a list of all available game server images
func (hr *httpRequestAuthenticator) ListImages(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.ListImages(c)
}

// GetStatus - Query status of all deployments
func (hr *httpRequestAuthenticator) GetStatus(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.GetStatus(c)
}

// ConfigureContainer - Configure a game server based on POST body
func (hr *httpRequestAuthenticator) ConfigureContainer(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.ConfigureContainer(c)
}

// DeployContainer - Deploy a game server based on POST body
func (hr *httpRequestAuthenticator) DeployContainer(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.DeployContainer(c)
}

// StartContainer - Start a game server/container
func (hr *httpRequestAuthenticator) StartContainer(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.StartContainer(c)
}

// StopContainer - Stop a game server/container
func (hr *httpRequestAuthenticator) StopContainer(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.StopContainer(c)
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestAuthenticator) RestartContainer(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.RestartContainer(c)
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestAuthenticator) DeleteContainer(c *gin.Context) {
	if false {
		return
	}
	hr.nextHandler.DeleteContainer(c)
}
