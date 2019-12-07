package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func start(c *gin.Context, api API) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if id2, err := strconv.ParseUint(id, 10, 64); err == nil {
		if result, err := api.Start(id2); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = result
			c.JSON(http.StatusAccepted, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}
}

func stop(c *gin.Context, api API) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if id2, err := strconv.ParseUint(id, 10, 64); err == nil {
		if result, err := api.Stop(id2); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = result
			c.JSON(http.StatusAccepted, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}
}

func restart(c *gin.Context, api API) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

	if id2, err := strconv.ParseUint(id, 10, 64); err == nil {
		if result, err := api.Restart(id2); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = result
			c.JSON(http.StatusAccepted, h)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
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

	if id2, err := strconv.ParseUint(id, 10, 64); err == nil {
		if err := api.Destroy(id2); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}

		go func() {
			time.Sleep(time.Second * time.Duration(10))
			api.RemoveDeployment(id2)
		}()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
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

	if id, err := api.Deploy(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		h := gin.H{"status": "ok"}
		h["id"] = id
		c.JSON(http.StatusCreated, h)
	}
}

func status(c *gin.Context, api API) {
	id, exists := c.GetQuery("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error"})
		return
	}

	if id2, err := strconv.ParseUint(id, 10, 64); err == nil {
		if result, err := api.Status(id2); err == nil {
			h := gin.H{"status": "ok"}
			h["message"] = result
			c.JSON(http.StatusOK, h)
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error"})
}
