package openapi

import (
	"context"
	"errors"
	"flag"
	"fmt"
	uuidGen "github.com/twinj/uuid"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

const defaultNamespace = "gamebaseprefix"
const defaultNamespaceUser = defaultNamespace + "-user-"

type kubernetesClient struct {
	Client *kubernetes.Clientset
}

func newKubernetesClientset() kubernetesClient {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset := kubernetes.NewForConfigOrDie(config)
	return kubernetesClient{clientset}
}

func (k kubernetesClient) GetGameServerList(ctx context.Context, namespace string) ([]*gameServer, error) {
	allDeployments, err := k.Client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	gameServers := make([]*gameServer, 0)
	for _, deployment := range allDeployments.Items {
		uuid, exists := deployment.Labels["deploymentUUID"]
		if !exists {
			fmt.Println("Found Deployment " + deployment.Name + " without deploymentUUID")
		}
		uuidGameServer, err := k.GetGameServer(ctx, namespace, uuid)
		if err != nil {
			fmt.Println("Could not find complete GameServer for UUID:" + uuid)
			fmt.Println(err)
		} else {
			gameServers = append(gameServers, uuidGameServer)
		}
	}
	return gameServers, nil
}

func (k kubernetesClient) GetGameServer(ctx context.Context, namespace string, uuid string) (*gameServer, error) {
	existingConfigMap, err := k.Client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingConfigMap.Items); num != 1 {
		return nil, errors.New("Number of selected ConfigMaps for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	existingPVC, err := k.Client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingPVC.Items); num != 1 {
		return nil, errors.New("Number of selected PVCs for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	existingDeployment, err := k.Client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingDeployment.Items); num != 1 {
		return nil, errors.New("Number of selected Deployments for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	existingService, err := k.Client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingService.Items); num != 1 {
		return nil, errors.New("Number of selected Services for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	return &gameServer{
		configmap:  kubernetesComponentConfigMap{existingConfigMap.Items[0]},
		pvc:        kubernetesComponentPVC{existingPVC.Items[0]},
		deployment: kubernetesComponentDeployment{existingDeployment.Items[0]},
		service:    kubernetesComponentService{existingService.Items[0]},
	}, nil
}

func (k kubernetesClient) DeployTemplate(ctx context.Context, namespace string, template *gameServerTemplate) (*gameServer, error) {
	deploymentPayload := template.GetUniqueGameServer()
	createOptions := metav1.CreateOptions{}
	//Deploy ConfigMap
	deployedConfigMap, err := k.Client.CoreV1().ConfigMaps(namespace).Create(ctx, &deploymentPayload.configmap.ConfigMap, createOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployed ConfigMap: " + deployedConfigMap.GetName())
	//Deploy PersistentVolumeClaim
	deployedPVC, err := k.Client.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, &deploymentPayload.pvc.PersistentVolumeClaim, createOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployed PVC: " + deployedPVC.GetName())
	//Adapt Deployment with name references to the generated Names
	payloadDeployment := deploymentPayload.deployment.Deployment
	//Replace Reference to ConfigMap
	for _, container := range payloadDeployment.Spec.Template.Spec.Containers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef.Name == "GameServerTemplateConfigMap" {
				envFrom.ConfigMapRef.Name = deployedConfigMap.GetName()
				fmt.Println("Referenced ConfigMap in Deployment: " + deployedConfigMap.GetName())
			}
		}
	}
	//Replace reference to PVC
	for _, volume := range payloadDeployment.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim.ClaimName == "GameServerTemplatePersistentVolumeClaim" {
			volume.PersistentVolumeClaim.ClaimName = deployedPVC.GetName()
			fmt.Println("Referenced PVC in Deployment: " + deployedPVC.GetName())
		}
	}
	//Deploy Deployment
	deployedDeployment, err := k.Client.AppsV1().Deployments(namespace).Create(ctx, &payloadDeployment, createOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployed Deployment: " + deployedDeployment.GetName())
	//Deploy Service
	deployedService, err := k.Client.CoreV1().Services(namespace).Create(ctx, &deploymentPayload.service.Service, createOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployed Service: " + deployedService.GetName())
	return &gameServer{
		configmap:  kubernetesComponentConfigMap{*deployedConfigMap},
		pvc:        kubernetesComponentPVC{*deployedPVC},
		deployment: kubernetesComponentDeployment{*deployedDeployment},
		service:    kubernetesComponentService{*deployedService},
	}, nil
}

func (k kubernetesClient) DeleteGameserver(ctx context.Context, namespace string, target *gameServer) error {
	deleteOptions := metav1.DeleteOptions{}
	err := k.Client.CoreV1().ConfigMaps(namespace).Delete(ctx, target.configmap.Name, deleteOptions)
	if err != nil {
		return err
	}
	err = k.Client.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, target.pvc.Name, deleteOptions)
	if err != nil {
		return err
	}
	err = k.Client.AppsV1().Deployments(namespace).Delete(ctx, target.deployment.Name, deleteOptions)
	if err != nil {
		return err
	}
	err = k.Client.CoreV1().Services(namespace).Delete(ctx, target.service.Name, deleteOptions)
	if err != nil {
		return err
	}
	return nil
}

func (k kubernetesClient) UpdateDeployedGameserver(ctx context.Context, namespace string, target *gameServer) (*gameServer, error) {
	updateOptions := metav1.UpdateOptions{}
	//Update ConfigMap
	updatedConfigMap, err := k.Client.CoreV1().ConfigMaps(namespace).Update(ctx, &target.configmap.ConfigMap, updateOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Update ConfigMap: " + updatedConfigMap.GetName())
	//Update PersistentVolumeClaim
	updatedPVC, err := k.Client.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, &target.pvc.PersistentVolumeClaim, updateOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Update PVC: " + updatedPVC.GetName())
	//Update Deployment
	updatedDeployment, err := k.Client.AppsV1().Deployments(namespace).Update(ctx, &target.deployment.Deployment, updateOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Update Deployment: " + updatedDeployment.GetName())
	//Deploy Service
	updatedService, err := k.Client.CoreV1().Services(namespace).Update(ctx, &target.service.Service, updateOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Update Service: " + updatedService.GetName())
	return &gameServer{
		configmap:  kubernetesComponentConfigMap{*updatedConfigMap},
		pvc:        kubernetesComponentPVC{*updatedPVC},
		deployment: kubernetesComponentDeployment{*updatedDeployment},
		service:    kubernetesComponentService{*updatedService},
	}, nil
}

func (k kubernetesClient) Rescale(ctx context.Context, namespace string, target *gameServer, targetReplicas int32) error {
	scale, err := k.Client.AppsV1().Deployments(namespace).GetScale(ctx, target.deployment.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = targetReplicas
	_, err = k.Client.AppsV1().Deployments(namespace).UpdateScale(ctx, target.deployment.Name, scale, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (k kubernetesClient) CreateDockerConfigSecret(ctx context.Context, namespace string, name string, base64secret string) (*v1.Secret, error) {
	//base64secret = "{\"auths\": {\"url.to.server\": {\"auth\": \"base64=\"}}}"
	secretMap := map[string]string{".dockerconfigjson": base64secret}
	return k.CreateSecret(ctx, namespace, name, v1.SecretTypeDockerConfigJson, secretMap)
}

// Create a kubernetes secret which stores the user information
func (k kubernetesClient) CreateUserSecret(ctx context.Context, namespace string, name string, user GamebaseUser) (*v1.Secret, error) {
	return k.CreateSecret(ctx, namespace, name, v1.SecretTypeOpaque, user.ToSecretData())
}

func (k kubernetesClient) CreateSecret(ctx context.Context, namespace string, name string, secretType v1.SecretType, stringData map[string]string) (*v1.Secret, error) {
	_, _ = k.CreateNamespace(ctx, namespace)
	secret := v1.Secret{
		Type: secretType,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: stringData,
	}
	createdSecret, err := k.Client.CoreV1().Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return createdSecret, nil
}

func (k kubernetesClient) SetDefaultServiceAccountImagePullSecret(ctx context.Context, namespace string, name string) (*v1.ServiceAccount, error) {
	defaultServiceAccount, err := k.Client.CoreV1().ServiceAccounts(namespace).Get(ctx, "default", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	defaultServiceAccount.ImagePullSecrets = append(defaultServiceAccount.ImagePullSecrets, v1.LocalObjectReference{Name: name})
	updatedServiceAccount, err := k.Client.CoreV1().ServiceAccounts(namespace).Update(ctx, defaultServiceAccount, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return updatedServiceAccount, nil
}

func (k kubernetesClient) GetSecret(ctx context.Context, namespace string, name string) (*v1.Secret, error) {
	return k.Client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k kubernetesClient) UpdateSecret(ctx context.Context, namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return k.Client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
}

// Set the user information either by creating the kubernetes secret
// or if it already exists, updating it
func (k kubernetesClient) SetUserSecret(ctx context.Context, email string, user GamebaseUser) error {
	encoded := encodeEmail(email)

	secret, err := k.CreateSecret(ctx, defaultNamespace, encoded, v1.SecretTypeOpaque, map[string]string{})
	if err != nil && !strings.HasSuffix(err.Error(), "already exists") {
		return err
	}

	secret, err = k.GetSecret(ctx, defaultNamespace, encoded)
	if err != nil {
		return err
	}

	// new users might not have uuid so we need to generate one
	if secret.Data == nil {
		random := make([]byte, 4)
		rand.Read(random)
		id := uuidGen.NewV5(uuidGen.NameSpaceURL, "game-base.de/backend/user", email, time.Now().UTC(), random)
		uuid := uuidGen.Formatter(id, uuidGen.FormatHex)
		secret.Data = map[string][]byte{
			"uuid": []byte(uuid),
		}
	}

	for key, value := range user.ToSecretData() {
		if value != "" {
			secret.Data[key] = []byte(value)
		}
	}

	_, err = k.UpdateSecret(ctx, defaultNamespace, secret)
	return err
}

// Lookup the uuid of the user
func (k kubernetesClient) GetUuid(ctx context.Context, email string) (string, error) {
	encoded := encodeEmail(email)
	secret, err := k.GetSecret(ctx, defaultNamespace, encoded)
	if err != nil {
		return "", err
	}

	if uuid, exists := secret.Data["uuid"]; exists {
		return string(uuid), nil
	}

	return "", errors.New("secret does not contain uuid")
}

func (k kubernetesClient) CreateNamespace(ctx context.Context, name string) (*v1.Namespace, error) {
	namespace := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return k.Client.CoreV1().Namespaces().Create(ctx, &namespace, metav1.CreateOptions{})
}

func (k kubernetesClient) GetUserSecret(ctx context.Context, email string) (*GamebaseUser, error) {
	_, err := k.CreateNamespace(ctx, defaultNamespace)
	if err != nil && !strings.HasSuffix(err.Error(), "already exists") {
		return nil, err
	}

	secret, err := k.GetSecret(ctx, defaultNamespace, encodeEmail(email))
	if err != nil {
		return nil, err
	}

	user := NewGamebaseUserFromSecretData(email, secret.Data)
	return &user, nil
}

func (k kubernetesClient) DeleteUserSecret(ctx context.Context, email string) error {
	deleteOptions := metav1.DeleteOptions{}
	return k.Client.CoreV1().Secrets(defaultNamespace).Delete(ctx, encodeEmail(email), deleteOptions)
}
