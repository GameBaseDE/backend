package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func start(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusTeapot, gin.H{"status": "error"})
	}
}

func stop(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusTeapot, gin.H{"status": "error"})
	}
}

func restart(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusTeapot, gin.H{"status": "error"})
	}
}

func destroyByQuery(c *gin.Context) {
	id := c.Query("id")
	destroy(id, c)
}

func destroyByParam(c *gin.Context) {
	id := c.Param("id")
	destroy(id, c)
}

func destroy(id string, c *gin.Context) {
	if id != "" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusTeapot, gin.H{"status": "error"})
	}
}

func deploy(c *gin.Context) {
	var body DeployContainerRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": body})
}

func status(c *gin.Context) {
	id, exists := c.GetQuery("id")
	if exists {
		msg := QueryContainerRequest{id, "test", []uint16{1, 2}, 42}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": msg})
	} else {
		c.JSON(http.StatusTeapot, gin.H{"status": "error", "message": nil})
	}
}
