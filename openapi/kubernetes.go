package openapi

import (
	"errors"
	"flag"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/autoscaling/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch2 "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ptr "k8s.io/utils/pointer"
	"path/filepath"
)

func deploymentTemplate() appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.Int32Ptr(1),
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
							Name:  "test",
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
}

func initkube() *kubernetes.Clientset {
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
	return clientset
}

type API struct {
	Client            *kubernetes.Clientset
	StandardNamespace string
}

func NewAPI() API {
	client := initkube()
	api := API{Client: client, StandardNamespace: "test"}
	api.GetNamespaceOrDie()
	return api
}

func (api API) GetNamespace() (*apiv1.Namespace, error) {
	return api.GetNamespaceClient().Get(api.StandardNamespace, metav1.GetOptions{})
}

func (api API) GetNamespaceClient() coreV1.NamespaceInterface {
	return api.Client.CoreV1().Namespaces()
}

func (api API) GetDeploymentClient() appsV1.DeploymentInterface {
	return api.Client.AppsV1().Deployments(api.StandardNamespace)
}

func (api API) GetNamespaceOrDie() apiv1.Namespace {
	if namespace, err := api.GetNamespace(); err != nil {
		template := apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: api.StandardNamespace,
			},
		}
		if _, err := api.GetNamespaceClient().Create(&template); err != nil {
			panic("Could not create standard namespace")
		}

		return template
	} else {
		return *namespace
	}
}

func (api API) Destroy(deploymentName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return api.GetDeploymentClient().Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (api API) Status(deploymentName string) (*appsv1.Deployment, error) {
	if deploymentName == "" {
		return nil, errors.New("Deployment with id " + deploymentName + " not found")
	}

	return api.GetDeploymentClient().Get(deploymentName, metav1.GetOptions{})
}

func (api API) Start(deploymentName string) (*v1.Scale, error) {
	return api.Rescale(deploymentName, 1)
}

func (api API) Stop(deploymentName string) (*v1.Scale, error) {
	return api.Rescale(deploymentName, 0)
}

func (api API) Restart(deploymentName string) (*appsv1.Deployment, error) {
	if deployment, err := api.Status(deploymentName); err != nil {
		return nil, err
	} else {
		if deployment.Status.Replicas == 0 {
			goto StartAgain
		}
	}

	if _, err := api.Stop(deploymentName); err != nil {
		return nil, err
	}

	if watch, err := api.GetDeploymentClient().Watch(metav1.ListOptions{Watch: true}); err == nil {
		defer watch.Stop()
		for event := range watch.ResultChan() {
			if event.Type != watch2.Error {
				deployment := event.Object.(*appsv1.Deployment)
				if deployment.Name == deploymentName {
					if deployment, err := api.Status(deploymentName); err != nil {
						return nil, err
					} else {
						if deployment.Status.Replicas == 0 {
							goto StartAgain
						}
					}
				}
			}
		}
	}

StartAgain:
	if _, err := api.Start(deploymentName); err != nil {
		return nil, err
	} else {
		return api.Status(deploymentName)
	}
}

func (api API) Rescale(deploymentName string, replicas int32) (*v1.Scale, error) {
	if deploymentName == "" {
		return nil, errors.New("Deployment with id " + deploymentName + " not found")
	}

	if scale, err := api.GetDeploymentClient().GetScale(deploymentName, metav1.GetOptions{}); err == nil && scale != nil {
		scale.Spec.Replicas = replicas
		return api.GetDeploymentClient().UpdateScale(deploymentName, scale)
	} else {
		return scale, err
	}
}

func (api API) Deploy(request GameServerConfigurationTemplate) (*appsv1.Deployment, error) {
	deployment := deploymentTemplate()
	deployment.Name = request.Id
	container := &deployment.Spec.Template.Spec.Containers[0]
	container.Image = request.Image
	container.Ports = []apiv1.ContainerPort{}

	for port := range request.Ports {
		container.Ports = append(container.Ports, apiv1.ContainerPort{
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: int32(port),
		})
	}

	return api.GetDeploymentClient().Create(&deployment)
}

func (api API) Configure(request GameServerConfigurationTemplate) (*appsv1.Deployment, error) {
	deploymentName := request.Id

	if err := api.Destroy(deploymentName); err != nil {
		return nil, err
	}

	return api.Deploy(request)
}

func (api API) List() ([]appsv1.Deployment, error) {
	if result, err := api.GetDeploymentClient().List(metav1.ListOptions{}); err != nil {
		return []appsv1.Deployment{}, err
	} else {
		return result.Items, nil
	}
}
