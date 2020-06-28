package openapi

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

type kubernetesComponentPVC struct {
	v1.PersistentVolumeClaim
}

func (k *kubernetesComponentPVC) Rename(templateName string) error {
	k.Name = ""
	k.GenerateName = templateName + "-pvc-"
	return nil
}

func (k *kubernetesComponentPVC) Validate(templatePrefix string) error {
	err := validateComponentName("PersistentVolumeClaim", templatePrefix, k.GenerateName)
	if err != nil {
		return err
	}
	err = checkBasicGameserverLabels(k.Labels)
	if err != nil {
		return err
	}
	fmt.Println("PersistentVolumeClaim validated!")
	return nil
}
