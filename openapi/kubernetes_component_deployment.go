package openapi

import (
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
)

type kubernetesComponentDeployment struct {
	appsv1.Deployment
}

func (k *kubernetesComponentDeployment) Rename(templateName string) error {
	k.Name = ""
	k.GenerateName = templateName + "-deployment-"
	return nil
}

func (k *kubernetesComponentDeployment) Validate(templatePrefix string) error {
	err := validateComponentName("Deployment", templatePrefix, k.GenerateName)
	if err != nil {
		return err
	}
	err = checkBasicGameserverLabels(k.Labels)
	if err != nil {
		return err
	}
	err = checkLabelValue(k.Labels, "name", "GameServerTemplateName")
	if err != nil {
		return err
	}
	//Check Deployment Template Contains Containers
	if len(k.Spec.Template.Spec.Containers) == 0 {
		return errors.New("Deployment.Spec.Template.Spec.Containers does not contain Containers!")
	}
	fmt.Println("Deployment validated!")
	return nil
}
