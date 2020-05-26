package openapi

import (
	"github.com/gin-gonic/gin"
)

var api = NewAPI()
var authenticator = newHttpRequestAuthenticator()

func ConfigureContainer_(c *gin.Context) {
	authenticator.ConfigureContainer(c)
}

func DeleteContainer_(c *gin.Context) {
	authenticator.DeleteContainer(c)
}

func DeployContainer_(c *gin.Context) {
	authenticator.DeployContainer(c)
}

func GetStatus_(c *gin.Context) {
	authenticator.GetStatus(c)
}

func ListImages_(c *gin.Context) {
	authenticator.ListImages(c)
}

func RestartContainer_(c *gin.Context) {
	authenticator.RestartContainer(c)
}

func StartContainer_(c *gin.Context) {
	authenticator.StartContainer(c)
}

func StopContainer_(c *gin.Context) {
	authenticator.StopContainer(c)
}
