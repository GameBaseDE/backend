package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	ptr "k8s.io/utils/pointer"
	"net/http"
	"os"
	"path/filepath"
)

type DeployContainerRequest struct {
	Image string   `json:"image"`
	Ports []uint16 `json:"ports"`
	Slots uint16   `json:"slots"`
}

type QueryContainerRequest struct {
	Id    string   `json:"id"`
	Image string   `json:"image"`
	Ports []uint16 `json:"ports"`
	Slots uint16   `json:"slots"`
}

func main() {
	initkube()
	router := gin.Default()

	router.GET("/api", status)
	router.GET("/api/start/:id", start)
	router.GET("/api/stop/:id", stop)
	router.GET("/api/restart/:id", restart)
	router.POST("/api/deploy", deploy)
	router.DELETE("/api/destroy", destroyByQuery)
	router.DELETE("/api/destroy/:id", destroyByParam)

	// Listen and serve on 0.0.0.0:8080
	if router.Run(":8080") != nil {
		println("Could not start the server")
	}
}

func start(c *gin.Context) {
	id := c.Param("id")
	createDep("demo-deploymente" + id)
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

var deploymentsClient v1.DeploymentInterface
var deployment *appsv1.Deployment

func initkube() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
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
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient = clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
}

func createDep(deployName string) {
	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deployName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	// Create Deployment
	fmt.Println("Creating deployment...")
	_, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	print("abc")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
