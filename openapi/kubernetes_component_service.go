package openapi

import (
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
)

type kubernetesComponentService struct {
	v1.Service
}

func (k *kubernetesComponentService) Rename(templateName string) error {
	k.Name = ""
	k.GenerateName = templateName + "-service-"
	return nil
}

func (k *kubernetesComponentService) Validate(templatePrefix string) error {
	err := validateComponentName("Service", templatePrefix, k.GenerateName)
	if err != nil {
		return err
	}
	err = checkBasicGameserverLabels(k.Labels)
	if err != nil {
		return err
	}
	//Check exposed Ports
	if len(k.Spec.Ports) == 0 {
		return errors.New("Service.Spec.Ports is empty!")
	}
	//Check Selector
	if len(k.Spec.Selector) == 0 {
		return errors.New("Service.Spec.Selector is unset!")
	}
	fmt.Println("Service validated!")
	return nil
}
