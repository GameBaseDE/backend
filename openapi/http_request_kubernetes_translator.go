package openapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpRequestKubernetesTranslator struct {
	nextHandler httpRequestHandler
	api         API
	cl          kubernetesClient
	templates   []*gameServerTemplate
}

func newHttpRequestKubernetesTranslator() *httpRequestKubernetesTranslator {
	return &httpRequestKubernetesTranslator{api: API{}, cl: newKubernetesClientset(), templates: readGameServerTemplates()}
}

// Login - Login a user and return a JWT with the user object
func (hr *httpRequestKubernetesTranslator) Login(c *gin.Context) {
	return
}

// Logout - Invalidate the passed JWT
func (hr *httpRequestKubernetesTranslator) Logout(c *gin.Context) {
	return
}

// Logout - Invalidate the passed JWT
func (hr *httpRequestKubernetesTranslator) Register(c *gin.Context) {
	return
}

// ListTemplates - Get a list of all available game server images
func (hr *httpRequestKubernetesTranslator) ListTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// GetStatus - Query status of all deployments
func (hr *httpRequestKubernetesTranslator) GetStatus(c *gin.Context) {
	id := c.GetString("id")
	if id == "" {
		if result, err := hr.api.List(); err == nil {
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

	if result, err := hr.api.Status(id); err == nil {
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
				if result, err := hr.api.Configure(id, request); err != nil {
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
			_ = request
			if err := hr.cl.DeployTemplate(c.GetString("namespace"), hr.templates[0]); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				h := gin.H{"status": "ok"}
				//h["message"] = AsGameServerStatus(result)
				h["message"] = GameContainerStatus{
					Id:                "1",
					Status:            "??",
					Configuration:     GameContainerConfiguration{},
					GameServerDetails: nil,
				}
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
	hr.cl.GetNamespaceList()
	hr.cl.GetConfigMaps(c.GetString("namespace"))
	hr.cl.GetPVCs(c.GetString("namespace"))
	hr.cl.GetDeploymentList(c.GetString("namespace"))
	hr.cl.GetServices(c.GetString("namespace"))
	hr.cl.GetPVCs2(c.GetString("namespace"), "gameserver=cuberite,deployment=gameserver")
	//hr.cl.UpdatePVC("gambaseprefix-testuser","cuberite")
	hr.cl.CreateDockerConfigSecret(c.GetString("namespace"), "regcred", "aaabbbcccddd")
	hr.cl.SetDefaultServiceAccountImagePullSecret(c.GetString("namespace"), "regcred")
	if id := c.GetString("id"); id != "" {
		if result, err := hr.api.Start(id); err == nil {
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
		if result, err := hr.api.Stop(id); err == nil {
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
		if result, err := hr.api.Restart(id); err == nil {
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
		if err := hr.api.Destroy(id); err != nil {
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
