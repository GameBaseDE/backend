/*
 * GameBase Communication API
 *
 * This is the REST API used as an communication layer between backend and frontend.
 *
 * API version: 2.0.0
 * Contact: gamebase@gahr.dev
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"github.com/gin-gonic/gin"
)

// ConfigureContainer - Configure a game server based on POST body
func ConfigureContainer(c *gin.Context) {
	authenticator.ConfigureContainer(c)
}

// DeleteContainer - Delete deployment of game server
func DeleteContainer(c *gin.Context) {
	authenticator.DeleteContainer(c)
}

// DeployContainer - Deploy a game server based on POST body
func DeployContainer(c *gin.Context) {
	authenticator.DeployContainer(c)
}

// GetStatus - Query status of all deployments
func GetStatus(c *gin.Context) {
	authenticator.GetStatus(c)
}

// ListTemplates - Get a list of all available game server templates
func ListTemplates(c *gin.Context) {
	authenticator.ListTemplates(c)
}

// RestartContainer - Restart a game server/container
func RestartContainer(c *gin.Context) {
	authenticator.RestartContainer(c)
}

// StartContainer - Start a game server/container
func StartContainer(c *gin.Context) {
	authenticator.StartContainer(c)
}

// StopContainer - Stop a game server/container
func StopContainer(c *gin.Context) {
	authenticator.StopContainer(c)
}