package openapi

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type API struct {
	Client            *kubernetes.Clientset
	StandardNamespace string
}

func NewAPI() API {
	client := initkube()
	api := API{Client: client, StandardNamespace: "test"}
	api.GetNamespaceOrDie()
	return api
}

func (api API) GetNamespace() (*apiv1.Namespace, error) {
	return api.GetNamespaceClient().Get(api.StandardNamespace, metav1.GetOptions{})
}

func (api API) GetNamespaceClient() coreV1.NamespaceInterface {
	return api.Client.CoreV1().Namespaces()
}

func (api API) GetDeploymentClient() appsV1.DeploymentInterface {
	return api.Client.AppsV1().Deployments(api.StandardNamespace)
}

func (api API) GetNamespaceOrDie() apiv1.Namespace {
	if namespace, err := api.GetNamespace(); err != nil {
		template := apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: api.StandardNamespace,
			},
		}
		if _, err := api.GetNamespaceClient().Create(&template); err != nil {
			panic("Could not create standard namespace")
		}

		return template
	} else {
		return *namespace
	}
}
