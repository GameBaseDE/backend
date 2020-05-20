package openapi

import (
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
)

var api = NewAPI()

func AsGameServerStatus(deployment *appsv1.Deployment) *GameServerStatus {
	return &GameServerStatus{
		Id:     deployment.Name,
		Image:  deployment.Spec.Template.Spec.Containers[0].Image,
		Status: make([]string, deployment.Status.Replicas),
	}
}

func ConfigureContainer_(c *gin.Context) {
	// TODO: openapi is not finished yet
	c.JSON(http.StatusOK, gin.H{})
}

func DeleteContainer_(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		id = c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
			return
		}
	}

	if err := api.Destroy(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

func DeployContainer_(c *gin.Context) {
	var request GameServerDeployTemplate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result, err := api.Deploy(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		h := gin.H{"status": "ok"}
		h["message"] = AsGameServerStatus(result)
		c.JSON(http.StatusCreated, h)
	}
}

func GetStatus_(c *gin.Context) {
	id, _ := c.GetQuery("id")

	if id == "" {
		if result, err := api.List(); err == nil {
			h := make([]GameServerStatus, 0)
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

func ListImages_(c *gin.Context) {
	// TODO: not specified
	c.JSON(http.StatusOK, gin.H{})
}

func RestartContainer_(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if result, err := api.Restart(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = AsGameServerStatus(result)
		c.JSON(http.StatusAccepted, h)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}
}

func StartContainer_(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if result, err := api.Start(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = result
		c.JSON(http.StatusAccepted, h)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}
}

func StopContainer_(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if result, err := api.Stop(id); err == nil {
		h := gin.H{"status": "ok"}
		h["message"] = result
		c.JSON(http.StatusAccepted, h)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}
}
