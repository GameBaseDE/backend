package openapi

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	instantiated *HttpRequestProcessingChain
	once         sync.Once
)

// This is Type to define our Chain of Command Pattern
type HttpRequestProcessingChain struct {
	nextHandler httpRequestHandler
}

// This Method returns a Singleton.
func NewHttpRequestProcessingChain() *HttpRequestProcessingChain {
	once.Do(func() {
		instantiated = &HttpRequestProcessingChain{nextHandler: newHttpRequestAuthenticator()}
	})
	return instantiated
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) Login(c *gin.Context) {
	hr.nextHandler.Login(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) Logout(c *gin.Context) {
	hr.nextHandler.Logout(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) Register(c *gin.Context) {
	hr.nextHandler.Register(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) ListTemplates(c *gin.Context) {
	hr.nextHandler.ListTemplates(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) GetStatus(c *gin.Context) {
	hr.nextHandler.GetStatus(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) ConfigureContainer(c *gin.Context) {
	hr.nextHandler.ConfigureContainer(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) DeployContainer(c *gin.Context) {
	hr.nextHandler.DeployContainer(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) StartContainer(c *gin.Context) {
	hr.nextHandler.StartContainer(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) StopContainer(c *gin.Context) {
	hr.nextHandler.StopContainer(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) RestartContainer(c *gin.Context) {
	hr.nextHandler.RestartContainer(c)
}

// Pass Function Call to first httpRequestHandler
func (hr *HttpRequestProcessingChain) DeleteContainer(c *gin.Context) {
	hr.nextHandler.DeleteContainer(c)
}
