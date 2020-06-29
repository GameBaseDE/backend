package openapi

import (
	"errors"
)

type gameServerTemplate struct {
	templateName string
	configmap    kubernetesComponentConfigMap
	pvc          kubernetesComponentPVC
	deployment   kubernetesComponentDeployment
	service      kubernetesComponentService
}

func (gst *gameServerTemplate) GetName() string {
	return gst.templateName
}

func (gst *gameServerTemplate) Rename() error {
	err := gst.configmap.Rename(gst.templateName)
	if err != nil {
		return err
	}
	err = gst.pvc.Rename(gst.templateName)
	if err != nil {
		return err
	}
	err = gst.deployment.Rename(gst.templateName)
	if err != nil {
		return err
	}
	err = gst.service.Rename(gst.templateName)
	if err != nil {
		return err
	}
	return nil
}

func (gst *gameServerTemplate) Validate(templatePrefix string) error {
	err := gst.configmap.Validate(templatePrefix)
	if err != nil {
		return err
	}
	err = gst.pvc.Validate(templatePrefix)
	if err != nil {
		return err
	}
	err = gst.deployment.Validate(templatePrefix)
	if err != nil {
		return err
	}
	err = gst.service.Validate(templatePrefix)
	if err != nil {
		return err
	}
	return nil
}

func validateComponentName(component string, templatePrefix string, found string) error {
	expecting := templatePrefix + component
	if expecting != found {
		errMsg := "Error prasing " + component + "! Expecting: " + expecting + " Found: " + found
		return errors.New(errMsg)
	} else {
		return nil
	}
}

func checkLabelExists(containedLabels map[string]string, key string) error {
	//if containedLabels != nil {
	if _, exists := containedLabels[key]; containedLabels != nil && exists {
		return nil
	}
	//}
	return errors.New("Label[" + key + "] does not exist!")
}

func checkLabelValue(containedLabels map[string]string, key string, value string) error {
	if labelValue, exists := containedLabels[key]; containedLabels != nil && exists {
		if labelValue == value {
			return nil
		}
	}
	return errors.New("Label[" + key + "]!=" + value)
}

func checkBasicGameserverLabels(containedLabels map[string]string) error {
	err := checkLabelExists(containedLabels, "gameserver")
	if err != nil {
		return err
	}
	err = checkLabelValue(containedLabels, "deploymentType", "gameserver")
	if err != nil {
		return err
	}
	return nil
}
