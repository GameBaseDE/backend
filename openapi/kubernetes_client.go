package openapi

import (
	"errors"
	"flag"
	"fmt"
	uuidGen "github.com/twinj/uuid"
	appsv1 "k8s.io/api/apps/v1"
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

func (k kubernetesClient) GetGameServerList(namespace string) ([]*gameServer, error) {
	allDeployments, err := k.Client.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	gameServers := make([]*gameServer, 0)
	for _, deployment := range allDeployments.Items {
		uuid, exists := deployment.Labels["deploymentUUID"]
		if !exists {
			fmt.Println("Found Deployment " + deployment.Name + " without deploymentUUID")
		}
		uuidGameServer, err := k.GetGameServer(namespace, uuid)
		if err != nil {
			fmt.Println("Could not find complete GameServer for UUID:" + uuid)
			fmt.Println(err)
		} else {
			gameServers = append(gameServers, uuidGameServer)
		}
	}
	return gameServers, nil
}

func (k kubernetesClient) GetGameServer(namespace string, uuid string) (*gameServer, error) {
	existingConfigMap, err := k.Client.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingConfigMap.Items); num != 1 {
		return nil, errors.New("Number of selected ConfigMaps for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	existingPVC, err := k.Client.CoreV1().PersistentVolumeClaims(namespace).List(metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingPVC.Items); num != 1 {
		return nil, errors.New("Number of selected PVCs for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	existingDeployment, err := k.Client.AppsV1().Deployments(namespace).List(metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
	if err != nil {
		return nil, err
	}
	if num := len(existingDeployment.Items); num != 1 {
		return nil, errors.New("Number of selected Deployments for UUID " + uuid + " == " + fmt.Sprint(num) + " should be 1")
	}
	existingService, err := k.Client.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector: "deploymentUUID=" + uuid})
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

func (k kubernetesClient) DeployTemplate(namespace string, template *gameServerTemplate) (*gameServer, error) {
	deploymentPayload := template.GetUniqueGameServer()
	//Deploy ConfigMap
	deployedConfigMap, err := k.Client.CoreV1().ConfigMaps(namespace).Create(&deploymentPayload.configmap.ConfigMap)
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployed ConfigMap: " + deployedConfigMap.GetName())
	//Deploy PersistentVolumeClaim
	deployedPVC, err := k.Client.CoreV1().PersistentVolumeClaims(namespace).Create(&deploymentPayload.pvc.PersistentVolumeClaim)
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
	deployedDeployment, err := k.Client.AppsV1().Deployments(namespace).Create(&payloadDeployment)
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployed Deployment: " + deployedDeployment.GetName())
	//Deploy Service
	deployedService, err := k.Client.CoreV1().Services(namespace).Create(&deploymentPayload.service.Service)
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

func (k kubernetesClient) DeleteGameserver(namespace string, target *gameServer) error {
	deleteOptions := metav1.DeleteOptions{}
	err := k.Client.CoreV1().ConfigMaps(namespace).Delete(target.configmap.Name, &deleteOptions)
	if err != nil {
		return err
	}
	err = k.Client.CoreV1().PersistentVolumeClaims(namespace).Delete(target.pvc.Name, &deleteOptions)
	if err != nil {
		return err
	}
	err = k.Client.AppsV1().Deployments(namespace).Delete(target.deployment.Name, &deleteOptions)
	if err != nil {
		return err
	}
	err = k.Client.CoreV1().Services(namespace).Delete(target.service.Name, &deleteOptions)
	if err != nil {
		return err
	}
	return nil
}

func (k kubernetesClient) Rescale(namespace string, target *gameServer, targetReplicas int32) error {
	scale, err := k.Client.AppsV1().Deployments(namespace).GetScale(target.deployment.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = targetReplicas
	_, err = k.Client.AppsV1().Deployments(namespace).UpdateScale(target.deployment.Name, scale)
	if err != nil {
		return err
	}
	return nil
}

func (k kubernetesClient) CreateDockerConfigSecret(namespace string, name string, base64secret string) (*v1.Secret, error) {
	//base64secret = "{\"auths\": {\"url.to.server\": {\"auth\": \"base64=\"}}}"
	secretMap := map[string]string{".dockerconfigjson": base64secret}
	return k.CreateSecret(namespace, name, v1.SecretTypeDockerConfigJson, secretMap)
}

func (k kubernetesClient) CreateSecret(namespace string, name string, secretType v1.SecretType, stringData map[string]string) (*v1.Secret, error) {
	_, _ = k.CreateNamespace(namespace)
	secret := v1.Secret{
		Type: secretType,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: stringData,
	}
	createdSecret, err := k.Client.CoreV1().Secrets(namespace).Create(&secret)
	if err != nil {
		return nil, err
	}
	return createdSecret, nil
}

func (k kubernetesClient) SetDefaultServiceAccountImagePullSecret(namespace string, name string) (*v1.ServiceAccount, error) {
	defaultServiceAccount, err := k.Client.CoreV1().ServiceAccounts(namespace).Get("default", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	defaultServiceAccount.ImagePullSecrets = append(defaultServiceAccount.ImagePullSecrets, v1.LocalObjectReference{Name: name})
	updatedServiceAccount, err := k.Client.CoreV1().ServiceAccounts(namespace).Update(defaultServiceAccount)
	if err != nil {
		return nil, err
	}
	return updatedServiceAccount, nil
}

func (k kubernetesClient) GetConfigMaps(namespace string) (*v1.ConfigMapList, error) {
	return k.Client.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
}

func (k kubernetesClient) GetPVCs(namespace string) (*v1.PersistentVolumeClaimList, error) {
	return k.Client.CoreV1().PersistentVolumeClaims(namespace).List(metav1.ListOptions{})
}

func (k kubernetesClient) GetDeploymentList(namespace string) (*appsv1.DeploymentList, error) {
	return k.Client.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
}

func (k kubernetesClient) GetServices(namespace string) (*v1.ServiceList, error) {
	return k.Client.CoreV1().Services(namespace).List(metav1.ListOptions{})
}

func (k kubernetesClient) GetNamespaceList() (*v1.NamespaceList, error) {
	return k.Client.CoreV1().Namespaces().List(metav1.ListOptions{})
}

func (k kubernetesClient) GetPVCs2(namespace string, selector string) (*v1.PersistentVolumeClaimList, error) {
	return k.Client.CoreV1().PersistentVolumeClaims(namespace).List(metav1.ListOptions{LabelSelector: selector})
}

func (k kubernetesClient) UpdatePVC(namespace string, name string) (*v1.PersistentVolumeClaim, error) {
	pvc, _ := k.Client.CoreV1().PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})
	pvc.Labels["key2"] = "value"
	return k.Client.CoreV1().PersistentVolumeClaims(namespace).Update(pvc)
}

func (k kubernetesClient) GetSecrets(namespace string, selector string) (*v1.SecretList, error) {
	return k.Client.CoreV1().Secrets(namespace).List(metav1.ListOptions{LabelSelector: selector})
}

func (k kubernetesClient) GetSecret(namespace string, name string) (*v1.Secret, error) {
	return k.Client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
}

func (k kubernetesClient) UpdateSecret(namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return k.Client.CoreV1().Secrets(namespace).Update(secret)
}

// Set the user information either by creating the kubernetes secret
// or if it already exists, updating it
func (k kubernetesClient) SetUserSecret(email string, user GamebaseUser) error {
	encoded := encodeEmail(email)

	secret, err := k.CreateSecret(defaultNamespace, encoded, v1.SecretTypeOpaque, map[string]string{})
	if err != nil && !strings.HasSuffix(err.Error(), "already exists") {
		return err
	}

	secret, err = k.GetSecret(defaultNamespace, encoded)
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

	_, err = k.UpdateSecret(defaultNamespace, secret)
	return err
}

// Lookup the uuid of the user
func (k kubernetesClient) GetUuid(email string) (string, error) {
	encoded := encodeEmail(email)
	secret, err := k.GetSecret(defaultNamespace, encoded)
	if err != nil {
		return "", err
	}

	if uuid, exists := secret.Data["uuid"]; exists {
		return string(uuid), nil
	}

	return "", errors.New("secret does not contain uuid")
}

func (k kubernetesClient) CreateNamespace(name string) (*v1.Namespace, error) {
	namespace := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return k.Client.CoreV1().Namespaces().Create(&namespace)
}

func (k kubernetesClient) GetUserSecret(email string) (*GamebaseUser, error) {
	_, err := k.CreateNamespace(defaultNamespace)
	if err != nil && !strings.HasSuffix(err.Error(), "already exists") {
		return nil, err
	}

	secret, err := k.GetSecret(defaultNamespace, encodeEmail(email))
	if err != nil {
		return nil, err
	}

	user := NewGamebaseUserFromSecretData(email, secret.Data)
	return &user, nil
}

func (k kubernetesClient) DeleteUserSecret(email string) error {
	deleteOptions := metav1.DeleteOptions{}
	return k.Client.CoreV1().Secrets(defaultNamespace).Delete(encodeEmail(email), &deleteOptions)
}
