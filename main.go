package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type DeployContainerRequest struct {
	Image string   `json:"image"`
	Ports []uint16 `json:"ports"`
	Slots uint16   `json:"slots"`
}

func main() {
	deploymentsClient := initkube()
	api := API{Client: deploymentsClient}

	router := gin.Default()

	router.GET("/api", func(c *gin.Context) { status(c, api) })
	router.GET("/api/start/:id", func(c *gin.Context) { start(c, api) })
	router.GET("/api/stop/:id", func(c *gin.Context) { stop(c, api) })
	router.GET("/api/restart/:id", func(c *gin.Context) { restart(c, api) })
	router.POST("/api/deploy", func(c *gin.Context) { deploy(c, api) })
	router.DELETE("/api/destroy", func(c *gin.Context) { destroyByQuery(c, api) })
	router.DELETE("/api/destroy/:id", func(c *gin.Context) { destroyByParam(c, api) })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if router.Run(":"+port) != nil {
		println("Could not start the server")
	}
}

func initkube() v1.DeploymentInterface {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset := kubernetes.NewForConfigOrDie(config)
	return clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
}
