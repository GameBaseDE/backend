package main

import (
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	ptr "k8s.io/utils/pointer"
)

func deploymentTemplate() appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
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
}

type API struct {
	Client      v1.DeploymentInterface
	deployments map[string]string
}

func (api API) GetDeploymentName(id string) string {
	return api.deployments[id]
}

func (api API) Destroy(id string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return api.Client.Delete(api.GetDeploymentName(id), &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (api API) Status(id string) error {
	return nil
}

func (api API) Deploy(body DeployContainerRequest) (*appsv1.Deployment, error) {
	deployment := deploymentTemplate()
	deployment.Spec.Template.Spec.Containers[0].Image = body.Image

	val, err := api.Client.Create(&deployment)
	if err == nil {
		var data []byte
		_ = json.Unmarshal(data, body)
		api.deployments[body.Image] = string(data)
	}

	return val, err
}
