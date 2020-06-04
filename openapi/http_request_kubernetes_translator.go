package openapi

import (
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

func AsGameServerStatus(deployment *appsv1.Deployment) *GameContainerStatus {
	status := UNKNOWN
	conditionsLength := len(deployment.Status.Conditions)
	if conditionsLength != 0 {
		latestCondition := deployment.Status.Conditions[conditionsLength-1]
		switch latestCondition.Status {
		case v1.ConditionTrue:
			switch latestCondition.Type {
			case appsv1.DeploymentReplicaFailure:
				status = ERROR
			case appsv1.DeploymentAvailable:
				status = RUNNING
			case appsv1.DeploymentProgressing:
				status = RESTARTING
			}
		case v1.ConditionFalse:
			status = ERROR
		case v1.ConditionUnknown:
			status = UNKNOWN
		}

		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 0 {
			status = STOPPED
		}
	}

	return &GameContainerStatus{
		Id:     deployment.Name,
		Status: status,
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
	id := c.GetString("id")
	if id == "" {
		if result, err := api.List(); err == nil {
			h := make([]GameContainerStatus, 0)
			for _, status := range result {
				h = append(h, *AsGameServerStatus(&status))
			}

			c.JSON(http.StatusOK, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	if result, err := api.Status(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = AsGameServerStatus(result)
		c.JSON(http.StatusOK, h)
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// ConfigureContainer - Configure a game server based on POST body
func (hr *httpRequestKubernetesTranslator) ConfigureContainer(c *gin.Context) {
	if id := c.GetString("id"); id != "" {
		if request, exists := c.Get("request"); exists {
			if request, ok := request.(GameContainerConfiguration); ok {
				if result, err := api.Configure(id, request); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					h := gin.H{"status": "ok"}
					h["message"] = AsGameServerStatus(result)
					c.JSON(http.StatusOK, h)
				}

				return
			}

			panic("request is of invalid type")
		}

		panic("request is unset")
	}

	panic("id is unset")
}

// DeployContainer - Deploy a game server based on POST body
func (hr *httpRequestKubernetesTranslator) DeployContainer(c *gin.Context) {
	if request, exists := c.Get("request"); exists {
		if request, ok := request.(GameContainerDeployment); ok {
			if result, err := api.Deploy(request); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				h := gin.H{"status": "ok"}
				h["message"] = AsGameServerStatus(result)
				c.JSON(http.StatusCreated, h)
			}

			return
		}

		panic("request is of invalid type")
	}

	panic("request is unset")
}

// StartContainer - Start a game server/container
func (hr *httpRequestKubernetesTranslator) StartContainer(c *gin.Context) {
	if id := c.GetString("id"); id != "" {
		if result, err := api.Start(id); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = result
			c.JSON(http.StatusAccepted, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	}

	panic("id is unset")
}

// StopContainer - Stop a game server/container
func (hr *httpRequestKubernetesTranslator) StopContainer(c *gin.Context) {
	if id := c.GetString("id"); id != "" {
		if result, err := api.Stop(id); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = result
			c.JSON(http.StatusAccepted, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	}

	panic("id is unset")
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestKubernetesTranslator) RestartContainer(c *gin.Context) {
	if id := c.GetString("id"); id != "" {
		if result, err := api.Restart(id); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = AsGameServerStatus(result)
			c.JSON(http.StatusAccepted, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	}

	panic("id is unset")
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestKubernetesTranslator) DeleteContainer(c *gin.Context) {
	if id := c.GetString("id"); id != "" {
		if err := api.Destroy(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
	}

	panic("id is unset")
}

// Tests if a GameServer Id exists
func (hr *httpRequestKubernetesTranslator) existstGameServer(id string) {
}