package main

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/autoscaling/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch2 "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	ptr "k8s.io/utils/pointer"
	"math/rand"
	"strconv"
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

type API struct {
	Client            *kubernetes.Clientset
	Deployments       *map[uint64]string
	StandardNamespace string
}

func NewAPI(client *kubernetes.Clientset) API {
	deployments := make(map[uint64]string)
	api := API{Client: client, Deployments: &deployments, StandardNamespace: "test"}
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

func (api API) GetDeploymentName(id uint64) string {
	deployments := *api.Deployments
	return deployments[id]
}

func (api API) AddDeployment(uid string) uint64 {
	deployments := *api.Deployments

	id := rand.Uint64()
	for isDuplicate := true; isDuplicate; id = rand.Uint64() {
		for k := range deployments {
			if k == id {
				isDuplicate = true
				break
			}
		}

		isDuplicate = false
	}

	deployments[id] = uid
	return id
}

func (api API) RemoveDeployment(id uint64) {
	deployments := *api.Deployments
	delete(deployments, id)
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

func (api API) Destroy(id uint64) error {
	deploymentName := api.GetDeploymentName(id)
	deletePolicy := metav1.DeletePropagationForeground
	return api.GetDeploymentClient().Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (api API) Status(id uint64) (*appsv1.Deployment, error) {
	deploymentName := api.GetDeploymentName(id)

	if deploymentName == "" {
		return nil, errors.New("Deployment with id " + strconv.FormatUint(id, 10) + " not found")
	}

	return api.GetDeploymentClient().Get(deploymentName, metav1.GetOptions{})
}

func (api API) Start(id uint64) (*v1.Scale, error) {
	return api.Rescale(id, 1)
}

func (api API) Stop(id uint64) (*v1.Scale, error) {
	return api.Rescale(id, 0)
}

func (api API) Restart(id uint64) (*appsv1.Deployment, error) {
	deploymentName := api.GetDeploymentName(id)

	if deployment, err := api.Status(id); err != nil {
		return nil, err
	} else {
		if deployment.Status.Replicas == 0 {
			goto StartAgain
		}
	}

	if _, err := api.Stop(id); err != nil {
		return nil, err
	}

	if watch, err := api.GetDeploymentClient().Watch(metav1.ListOptions{Watch: true}); err == nil {
		defer watch.Stop()
		for event := range watch.ResultChan() {
			if event.Type != watch2.Error {
				deployment := event.Object.(*appsv1.Deployment)
				if deployment.Name == deploymentName {
					if deployment, err := api.Status(id); err != nil {
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
	if _, err := api.Start(id); err != nil {
		return nil, err
	} else {
		return api.Status(id)
	}
}

func (api API) Rescale(id uint64, replicas int32) (*v1.Scale, error) {
	deploymentName := api.GetDeploymentName(id)

	if deploymentName == "" {
		return nil, errors.New("Deployment with id " + strconv.FormatUint(id, 10) + " not found")
	}

	if scale, err := api.GetDeploymentClient().GetScale(deploymentName, metav1.GetOptions{}); err == nil && scale != nil {
		scale.Spec.Replicas = replicas
		return api.GetDeploymentClient().UpdateScale(deploymentName, scale)
	} else {
		return scale, err
	}
}

func (api API) Deploy(body DeployContainerRequest) (uint64, error) {
	deployment := deploymentTemplate()
	container := deployment.Spec.Template.Spec.Containers[0]
	container.Image = body.Image
	container.Ports = []apiv1.ContainerPort{}

	for port := range body.Ports {
		container.Ports = append(container.Ports, apiv1.ContainerPort{
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: int32(port),
		})
	}

	val, err := api.GetDeploymentClient().Create(&deployment)
	if err == nil {
		return api.AddDeployment(val.Name), nil
	}

	return 0, err
}
