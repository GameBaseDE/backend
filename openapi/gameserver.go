package openapi

type gameServer struct {
	configmap  kubernetesComponentConfigMap
	pvc        kubernetesComponentPVC
	deployment kubernetesComponentDeployment
	service    kubernetesComponentService
}
