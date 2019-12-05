package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func start(c *gin.Context, api API) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
	}
}

func stop(c *gin.Context, api API) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
	}
}

func restart(c *gin.Context, api API) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
	}
}

func destroyByQuery(c *gin.Context, api API) {
	id := c.Query("id")
	destroy(id, c, api)
}

func destroyByParam(c *gin.Context, api API) {
	id := c.Param("id")
	destroy(id, c, api)
}

func destroy(id string, c *gin.Context, api API) {
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if err := api.Destroy(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

func deploy(c *gin.Context, api API) {
	var body DeployContainerRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result, err := api.Deploy(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		h := gin.H{"status": "ok"}
		h["message"] = result
		c.JSON(http.StatusCreated, h)
	}
}

func status(c *gin.Context, api API) {
	id, exists := c.GetQuery("id")
	if exists {
		c.JSON(http.StatusAccepted, gin.H{"status": "ok", "message": id})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "error"})
	}
}
