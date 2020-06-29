package openapi

import (
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
)

type kubernetesComponentConfigMap struct {
	v1.ConfigMap
}

func (k *kubernetesComponentConfigMap) Rename(templateName string) error {
	k.Name = ""
	k.GenerateName = templateName + "-configmap-"
	return nil
}

func (k *kubernetesComponentConfigMap) Validate(templatePrefix string) error {
	err := validateComponentName("ConfigMap", templatePrefix, k.GenerateName)
	if err != nil {
		return err
	}
	err = checkBasicGameserverLabels(k.Labels)
	if err != nil {
		return err
	}
	//Check Data
	if len(k.Data) == 0 {
		return errors.New("ConfigMap.Data is unset!")
	}
	fmt.Println("ConfigMap validated!")
	return nil
}
