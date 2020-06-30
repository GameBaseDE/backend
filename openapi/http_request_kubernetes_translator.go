package openapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

func (hr *httpRequestKubernetesTranslator) kubernetesClient() kubernetesClient {
	return hr.cl
}

// Login - Login a user and return a JWT with the user object
func (hr *httpRequestKubernetesTranslator) Login(c *gin.Context) {
	return
}

// Logout - Invalidate the passed JWT
func (hr *httpRequestKubernetesTranslator) Logout(c *gin.Context) {
	return
}

// Register - Register a user and return a JWT with the user object
func (hr *httpRequestKubernetesTranslator) Register(c *gin.Context) {
	request, exists := c.Get("request")
	if !exists {
		panic("request is unset")
	}
	user, validUser := request.(GamebaseUser)
	if !validUser {
		panic("request is of invalid type")
	}

	uuid, _, err := hr.cl.GetUuid(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	namespace := defaultNamespaceUser + uuid
	if err := hr.cl.SetUserSecret(namespace, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, _, err := createToken(user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, User{
		Email:    user.Email,
		FullName: user.Name,
		Token:    token,
	})
	return
}

// ListTemplates - Get a list of all available game server images
func (hr *httpRequestKubernetesTranslator) ListTemplates(c *gin.Context) {
	if hr.templates == nil {
		c.JSON(http.StatusInternalServerError, Exception{Details: "templates not parsed!"})
		return
	}
	templatesList := []string{}
	for _, template := range hr.templates {
		templatesList = append(templatesList, template.GetName())
	}
	c.JSON(http.StatusOK, templatesList)
}

// GetStatus - Query status of all deployments
func (hr *httpRequestKubernetesTranslator) GetStatus(c *gin.Context) {
	id := c.GetString("id")
	existingGameServers := []*gameServer{}
	if id == "" {
		gameServers, err := hr.cl.GetGameServerList(getNamespace(c))
		existingGameServers = gameServers
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		_, existingGameServer := hr.parseIdRequest(c)
		if existingGameServer == nil {
			return
		}
		existingGameServers = append(existingGameServers, existingGameServer)
	}
	gameContainerStatuses := []*GameContainerStatus{}
	for _, gameServer := range existingGameServers {
		gameContainerStatus := gameServer.readGameContainerStatus()
		gameContainerStatuses = append(gameContainerStatuses, &gameContainerStatus)
	}
	c.JSON(http.StatusOK, gameContainerStatuses)
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
	request, exists := c.Get("request")
	if !exists {
		panic("request is unset")
	}
	deploymentRequest, ok := request.(GameContainerDeployment)
	if !ok {
		panic("request is of invalid type")
	}
	template, err := findGameServerTemplate(deploymentRequest.TemplatePath, hr.templates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	gameServer, err := hr.cl.DeployTemplate(getNamespace(c), template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Exception{Id: gameServer.GetUID(), Details: err.Error()})
		return
	}
	h := gin.H{"status": "ok"}
	c.JSON(http.StatusCreated, h)
	return
}

// StartContainer - Start a game server/container
func (hr *httpRequestKubernetesTranslator) StartContainer(c *gin.Context) {
	if hr.rescale(c, 1) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// StopContainer - Stop a game server/container
func (hr *httpRequestKubernetesTranslator) StopContainer(c *gin.Context) {
	if hr.rescale(c, 0) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestKubernetesTranslator) RestartContainer(c *gin.Context) {
	// Stop
	if !hr.rescale(c, 0) {
		return
	}
	// Wait until stopped
	_, existingGameServer := hr.parseIdRequest(c)
	maxruntime := *existingGameServer.GetTerminationTimeout()
	for i := int64(0); i <= maxruntime; i++ {
		if existingGameServer.GetStatus() == STOPPED {
			break
		}
		time.Sleep(1 * time.Second)
		_, existingGameServer = hr.parseIdRequest(c)
	}
	// Start
	if !hr.rescale(c, 1) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestKubernetesTranslator) DeleteContainer(c *gin.Context) {
	namespace, existingGameServer := hr.parseIdRequest(c)
	if existingGameServer == nil {
		return
	}
	err := hr.cl.DeleteGameserver(namespace, existingGameServer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Exception{Id: "", Details: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Tests if a GameServer Id exists
func (hr *httpRequestKubernetesTranslator) existstGameServer(id string) {
}

func getNamespace(c *gin.Context) string {
	if namespace := c.GetString("namespace"); namespace != "" {
		return namespace
	} else {
		panic("No Namespace in gin Context")
	}
}

// This method is used to parse all requests that specify the target Gameserver in the URL
func (hr *httpRequestKubernetesTranslator) parseIdRequest(c *gin.Context) (string, *gameServer) {
	id := c.GetString("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, Exception{Id: "", Details: "No ID specified"})
		return "", nil
	}
	namespace := getNamespace(c)
	existingGameServer, err := hr.cl.GetGameServer(namespace, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Exception{Id: "", Details: err.Error()})
		return "", nil
	}
	return namespace, existingGameServer
}

// rescale is used for Start,Stop and Restart
func (hr *httpRequestKubernetesTranslator) rescale(c *gin.Context, replicas int32) bool {
	namespace, existingGameServer := hr.parseIdRequest(c)
	if existingGameServer == nil {
		return false
	}
	err := hr.cl.Rescale(namespace, existingGameServer, replicas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Exception{Id: existingGameServer.GetUID(), Details: err.Error()})
		return false
	}
	return true
}
