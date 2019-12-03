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

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if router.Run(":"+port) != nil {
		println("Could not start the server")
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
