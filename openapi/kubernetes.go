package openapi

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/autoscaling/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch2 "k8s.io/apimachinery/pkg/watch"
)

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

func (api API) Deploy(request GameContainerDeployment) (*appsv1.Deployment, error) {
	deployment := deploymentTemplate()
	container := &deployment.Spec.Template.Spec.Containers[0]
	container.Image = request.TemplatePath

	return api.GetDeploymentClient().Create(&deployment)
}

func (api API) Configure(id string, request GameContainerConfiguration) (*appsv1.Deployment, error) {
	if err := api.Destroy(id); err != nil {
		return nil, err
	}

	deployment := deploymentTemplate()
	deployment.ObjectMeta.Name = request.Details.ServerName
	container := &deployment.Spec.Template.Spec.Containers[0]
	container.Image = request.Resources.TemplatePath
	container.Ports = []apiv1.ContainerPort{}

	return api.GetDeploymentClient().Create(&deployment)
}

func (api API) List() ([]appsv1.Deployment, error) {
	if result, err := api.GetDeploymentClient().List(metav1.ListOptions{}); err != nil {
		return []appsv1.Deployment{}, err
	} else {
		return result.Items, nil
	}
}
