package openapi

import (
	"errors"
	uuidGen "github.com/twinj/uuid"
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

func (gst *gameServerTemplate) GetUniqueGameServer() gameServer {
	uniqueGameServer := gameServer{
		configmap:  gst.configmap.DeeperCopy(),
		pvc:        gst.pvc.DeeperCopy(),
		deployment: gst.deployment.DeeperCopy(),
		service:    gst.service.DeeperCopy(),
	}
	//Populate user-definable Gameserver name
	uniqueGameServer.deployment.Labels["name"] = gst.templateName
	//Generate UUID
	uuid := uuidGen.NewV4().String()
	//Set UUID in metadata.Labels
	uniqueGameServer.configmap.Labels["deploymentUUID"] = uuid
	uniqueGameServer.pvc.Labels["deploymentUUID"] = uuid
	uniqueGameServer.deployment.Labels["deploymentUUID"] = uuid
	uniqueGameServer.deployment.Spec.Template.Labels["deploymentUUID"] = uuid
	uniqueGameServer.service.Labels["deploymentUUID"] = uuid
	//Get Label reference to deployment.spec.template.metadata.labels
	deploymentSelectorLabels := uniqueGameServer.deployment.Spec.Template.Labels
	//Adjust Selectors
	uniqueGameServer.deployment.Spec.Selector.MatchLabels = deploymentSelectorLabels
	uniqueGameServer.service.Spec.Selector = deploymentSelectorLabels
	return uniqueGameServer
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
