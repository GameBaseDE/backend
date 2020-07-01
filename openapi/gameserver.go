package openapi

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strings"
)

type gameServer struct {
	configmap  kubernetesComponentConfigMap
	pvc        kubernetesComponentPVC
	deployment kubernetesComponentDeployment
	service    kubernetesComponentService
}

func (gs *gameServer) getContainer() v1.Container {
	return gs.deployment.Spec.Template.Spec.Containers[0]
}

func (gs *gameServer) GetUID() string {
	return gs.deployment.Labels["deploymentUUID"]
}

func (gs *gameServer) GetStatus() Status {
	//gs.deployment.Status
	if gs.deployment.Status.Replicas == 0 {
		return STOPPED
	} else if gs.deployment.Status.UnavailableReplicas > 0 {
		return STARTING
	} else if gs.deployment.Status.AvailableReplicas > 0 {
		return RUNNING
	} else {
		return UNKNOWN
	}
}

func (gs *gameServer) GetName() string {
	return gs.deployment.Labels["name"]
}

func (gs *gameServer) SetName(newName string) {
	if newName == "" {
		return
	}
	gs.deployment.Labels["name"] = newName
}

func (gs *gameServer) GetTemplate() string {
	return gs.deployment.Labels["gameserver"]
}

func (gs *gameServer) GetPortMapping() []PortMapping {
	portMappings := []PortMapping{}
	for _, servicePort := range gs.service.Spec.Ports {
		if servicePort.Protocol == v1.ProtocolTCP {
			portMappings = append(portMappings, PortMapping{
				Protocol:      TCP,
				NodePort:      servicePort.NodePort,
				ContainerPort: servicePort.Port})
		} else if servicePort.Protocol == v1.ProtocolUDP {
			portMappings = append(portMappings, PortMapping{
				Protocol:      UDP,
				NodePort:      servicePort.NodePort,
				ContainerPort: servicePort.Port})
		}
	}
	return portMappings
}

func (gs *gameServer) SetPortMapping(newPortMapping []PortMapping) {
	if newPortMapping == nil {
		return
	}
	servicePorts := []v1.ServicePort{}
	for num, portMapping := range newPortMapping {
		newProtocol := v1.ProtocolTCP
		if portMapping.Protocol == UDP {
			newProtocol = v1.ProtocolUDP
		}
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:     "frontend-manual-" + fmt.Sprint(num),
			Protocol: newProtocol,
			Port:     portMapping.ContainerPort,
		})
	}
	gs.service.Spec.Ports = servicePorts
}

func (gs *gameServer) GetContainerMemoryLimit() int32 {
	//FIXME pass trough as string
	limits := gs.getContainer().Resources.Limits
	intLimit, converted := limits.Memory().AsInt64()
	if !converted {
		return -1
	}
	//TODO update openapi to int64
	return int32(intLimit) //What could possibly go wrong?
}

// SetContainerMemoryLimit - Recreates the PodSpec since Limits are immutable on existing objects
func (gs *gameServer) SetContainerMemoryLimit(newLimit int32) {
	if newLimit <= 0 {
		return
	}
	newMemoryLimitQuantity, err := resource.ParseQuantity(fmt.Sprint(newLimit))
	if err != nil {
		fmt.Println("Could not parse requested MemoryLimit")
		return
	}
	var newContainerSlice []v1.Container
	for _, container := range gs.deployment.Spec.Template.Spec.Containers {
		newResourceList := v1.ResourceList{}
		for resourceName, resourceQuantity := range container.Resources.Limits {
			if resourceName != v1.ResourceMemory {
				newResourceList[resourceName] = resourceQuantity
			} else {
				newResourceList[resourceName] = newMemoryLimitQuantity
			}
		}
		newResourceRequirements := v1.ResourceRequirements{
			Limits:   newResourceList,
			Requests: nil,
		}
		newContainer := v1.Container{
			Name:                     container.Name,
			Image:                    container.Image,
			Command:                  container.Command,
			Args:                     container.Args,
			WorkingDir:               container.WorkingDir,
			Ports:                    container.Ports,
			EnvFrom:                  container.EnvFrom,
			Env:                      container.Env,
			Resources:                newResourceRequirements,
			VolumeMounts:             container.VolumeMounts,
			VolumeDevices:            container.VolumeDevices,
			LivenessProbe:            container.LivenessProbe,
			ReadinessProbe:           container.ReadinessProbe,
			StartupProbe:             container.StartupProbe,
			Lifecycle:                container.Lifecycle,
			TerminationMessagePath:   container.TerminationMessagePath,
			TerminationMessagePolicy: container.TerminationMessagePolicy,
			ImagePullPolicy:          container.ImagePullPolicy,
			SecurityContext:          container.SecurityContext,
			Stdin:                    container.Stdin,
			StdinOnce:                container.StdinOnce,
			TTY:                      container.TTY,
		}
		newContainerSlice = append(newContainerSlice, newContainer)
	}
	podTemplateSpec := gs.deployment.Spec.Template.Spec
	newPodSpec := v1.PodSpec{
		Volumes:                       podTemplateSpec.Volumes,
		InitContainers:                podTemplateSpec.InitContainers,
		Containers:                    newContainerSlice,
		EphemeralContainers:           podTemplateSpec.EphemeralContainers,
		RestartPolicy:                 podTemplateSpec.RestartPolicy,
		TerminationGracePeriodSeconds: podTemplateSpec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         podTemplateSpec.ActiveDeadlineSeconds,
		DNSPolicy:                     podTemplateSpec.DNSPolicy,
		NodeSelector:                  podTemplateSpec.NodeSelector,
		ServiceAccountName:            podTemplateSpec.ServiceAccountName,
		AutomountServiceAccountToken:  podTemplateSpec.AutomountServiceAccountToken,
		NodeName:                      podTemplateSpec.NodeName,
		HostNetwork:                   podTemplateSpec.HostNetwork,
		HostPID:                       podTemplateSpec.HostPID,
		HostIPC:                       podTemplateSpec.HostIPC,
		ShareProcessNamespace:         podTemplateSpec.ShareProcessNamespace,
		SecurityContext:               podTemplateSpec.SecurityContext,
		ImagePullSecrets:              podTemplateSpec.ImagePullSecrets,
		Hostname:                      podTemplateSpec.Hostname,
		Subdomain:                     podTemplateSpec.Subdomain,
		Affinity:                      podTemplateSpec.Affinity,
		SchedulerName:                 podTemplateSpec.SchedulerName,
		Tolerations:                   podTemplateSpec.Tolerations,
		HostAliases:                   podTemplateSpec.HostAliases,
		PriorityClassName:             podTemplateSpec.PriorityClassName,
		Priority:                      podTemplateSpec.Priority,
		DNSConfig:                     podTemplateSpec.DNSConfig,
		ReadinessGates:                podTemplateSpec.ReadinessGates,
		RuntimeClassName:              podTemplateSpec.RuntimeClassName,
		EnableServiceLinks:            podTemplateSpec.EnableServiceLinks,
		PreemptionPolicy:              podTemplateSpec.PreemptionPolicy,
		Overhead:                      podTemplateSpec.Overhead,
		TopologySpreadConstraints:     podTemplateSpec.TopologySpreadConstraints,
	}
	gs.deployment.Spec.Template.Spec = newPodSpec
}

func (gs *gameServer) GetStartupArgs() string {
	return strings.Join(gs.getContainer().Args, " ")
}

func (gs *gameServer) SetStartupArgs(newArgs string) {
	if newArgs == "" {
		return
	}
	container := gs.getContainer()
	container.Args = strings.Split(newArgs, " ")
}

func (gs *gameServer) GetDescription() string {
	//TODO Evaluate using a separate ConfigMap or encoding into Labels
	return gs.configmap.Data["DESCRIPTION"]
}

func (gs *gameServer) SetDescription(newDescription string) {
	if newDescription == "" {
		return
	}
	//TODO Evaluate using a separate ConfigMap or encoding into Labels
	gs.configmap.Data["DESCRIPTION"] = newDescription
}

func (gs *gameServer) GetRestartBehavior() RestartBehavior {
	switch gs.deployment.Spec.Template.Spec.RestartPolicy {
	case v1.RestartPolicyAlways:
		return ALWAYS
	case v1.RestartPolicyOnFailure:
		return ON_FAILURE
	case v1.RestartPolicyNever:
		return NONE
	}
	return NONE
}

func (gs *gameServer) SetRestartBehavior(newRestartBehavior RestartBehavior) {
	switch newRestartBehavior {
	case NONE:
		gs.deployment.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyNever
	case UNLESS_STOPPED:
		gs.deployment.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
	case ON_FAILURE:
		gs.deployment.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
	case ALWAYS:
		gs.deployment.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyAlways
	}
}

func (gs *gameServer) GetEnvironmentVars() map[string]string {
	return gs.configmap.Data
}

func (gs *gameServer) SetEnvironmentVars(newEnvs map[string]string) {
	if len(newEnvs) == 0 {
		return
	}
	gs.configmap.Data = newEnvs
}

func (gs *gameServer) readGameContainerStatus() GameContainerStatus {
	gsst := GameContainerStatus{
		Id:     gs.GetUID(),
		Status: gs.GetStatus(),
		Configuration: GameContainerConfiguration{
			Details: GameContainerConfigurationDetails{
				ServerName:  gs.GetName(),
				Description: gs.GetDescription(),
			},
			Resources: GameContainerConfigurationResources{
				TemplatePath:    gs.GetTemplate(),
				Ports:           gs.GetPortMapping(),
				Memory:          gs.GetContainerMemoryLimit(),
				StartupArgs:     gs.GetStartupArgs(),
				RestartBehavior: gs.GetRestartBehavior(),
				EnvironmentVars: gs.GetEnvironmentVars(),
			},
		},
		GameServerDetails: map[string]string{"Details": "None"},
	}
	return gsst
}

func (gs *gameServer) GetTerminationTimeout() *int64 {
	return gs.deployment.Spec.Template.Spec.TerminationGracePeriodSeconds
}

func (gs *gameServer) UpdateGameServer(configurationUpdate GameContainerConfiguration) (*gameServer, error) {
	updatedGameServer := gameServer{
		configmap:  gs.configmap,
		pvc:        gs.pvc,
		deployment: gs.deployment,
		service:    gs.service,
	}
	updatedGameServer.SetName(configurationUpdate.Details.ServerName)
	updatedGameServer.SetDescription(configurationUpdate.Details.Description)
	updatedGameServer.SetPortMapping(configurationUpdate.Resources.Ports)
	updatedGameServer.SetContainerMemoryLimit(configurationUpdate.Resources.Memory)
	updatedGameServer.SetStartupArgs(configurationUpdate.Resources.StartupArgs)
	updatedGameServer.SetRestartBehavior(configurationUpdate.Resources.RestartBehavior)
	updatedGameServer.SetEnvironmentVars(configurationUpdate.Resources.EnvironmentVars)
	return &updatedGameServer, nil
}
