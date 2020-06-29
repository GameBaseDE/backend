package openapi

import (
	v1 "k8s.io/api/core/v1"
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
	return string(gs.deployment.UID)
}

func (gs *gameServer) GetStatus() Status {
	//TODO implement
	return UNKNOWN
}

func (gs *gameServer) GetName() string {
	return gs.deployment.Labels["name"]
}

func (gs *gameServer) GetTemplate() string {
	return gs.deployment.Labels["gameserver"]
}

func (gs *gameServer) GetContainerPorts() GameContainerConfigurationResourcesPorts {
	tcp := []int32{}
	udp := []int32{}
	for _, servicePort := range gs.service.Spec.Ports {
		if servicePort.Protocol == v1.ProtocolTCP {
			tcp = append(tcp, servicePort.Port)
		} else if servicePort.Protocol == v1.ProtocolUDP {
			udp = append(udp, servicePort.Port)
		}

	}
	return GameContainerConfigurationResourcesPorts{
		Tcp: tcp,
		Udp: udp,
	}
}

func (gs *gameServer) GetContainerMemoryLimit() int32 {
	//FIXME pass trough as string
	limits := gs.getContainer().Resources.Limits
	intLimit, converted := limits.Memory().AsInt64()
	if !converted {
		return -1
	}
	return int32(intLimit) //What could possibly go wrong?
}

func (gs *gameServer) GetStartupArgs() string {
	return strings.Join(gs.getContainer().Args, " ")
}

func (gs *gameServer) GetDescription() string {
	//TODO introduce description label
	return ""
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

func (gs *gameServer) GetEnvironmentVars() map[string]string {
	return gs.configmap.Data
}

func (gs *gameServer) readGameContainerStatus() GameContainerStatus {
	gsst := GameContainerStatus{
		Id:     gs.GetUID(),
		Status: gs.GetStatus(),
		Configuration: GameContainerConfiguration{
			Details: GameContainerConfigurationDetails{
				ServerName:  gs.GetName(),
				OwnerId:     "", //FIXME remove
				Description: gs.GetDescription(),
			},
			Resources: GameContainerConfigurationResources{
				TemplatePath:    gs.GetTemplate(),
				Ports:           gs.GetContainerPorts(),
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