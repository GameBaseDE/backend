package openapi

import (
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
)

func AsGameServerStatus(deployment *appsv1.Deployment) *GameContainerStatus {
	return &GameContainerStatus{
		Id: deployment.Name,
	}
}

type httpRequestKubernetesTranslator struct {
	nextHandler httpRequestHandler
}

func newHttpRequestKubernetesTranslator() *httpRequestKubernetesTranslator {
	return &httpRequestKubernetesTranslator{}
}

// ListTemplates - Get a list of all available game server images
func (hr *httpRequestKubernetesTranslator) ListTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// GetStatus - Query status of all deployments
func (hr *httpRequestKubernetesTranslator) GetStatus(c *gin.Context) {
	id, _ := c.GetQuery("id")

	if id == "" {
		if result, err := api.List(); err == nil {
			h := make([]GameContainerStatus, 0)
			for _, status := range result {
				h = append(h, *AsGameServerStatus(&status))
			}

			c.JSON(http.StatusOK, h)
			return
		}
	}

	if result, err := api.Status(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = AsGameServerStatus(result)
		c.JSON(http.StatusOK, h)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error"})
}

// ConfigureContainer - Configure a game server based on POST body
func (hr *httpRequestKubernetesTranslator) ConfigureContainer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// DeployContainer - Deploy a game server based on POST body
func (hr *httpRequestKubernetesTranslator) DeployContainer(c *gin.Context) {
	//TODO test
	var request GameContainerDeployment
	if result, err := api.Deploy(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		h := gin.H{"status": "ok"}
		h["message"] = AsGameServerStatus(result)
		c.JSON(http.StatusCreated, h)
	}
}

// StartContainer - Start a game server/container
func (hr *httpRequestKubernetesTranslator) StartContainer(c *gin.Context) {
	id := c.Param("id")
	if result, err := api.Start(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = result
		c.JSON(http.StatusAccepted, h)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}
}

// StopContainer - Stop a game server/container
func (hr *httpRequestKubernetesTranslator) StopContainer(c *gin.Context) {
	id := c.Param("id")
	if result, err := api.Stop(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = result
		c.JSON(http.StatusAccepted, h)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestKubernetesTranslator) RestartContainer(c *gin.Context) {
	id := c.Param("id")
	if result, err := api.Restart(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = AsGameServerStatus(result)
		c.JSON(http.StatusAccepted, h)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestKubernetesTranslator) DeleteContainer(c *gin.Context) {
	id := c.Param("id")
	if err := api.Destroy(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

// Tests if a GameServer Id exists
func (hr *httpRequestKubernetesTranslator) existstGameServer(id string) {
}
