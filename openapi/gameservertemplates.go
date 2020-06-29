package openapi

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	sch "k8s.io/client-go/kubernetes/scheme"
	"log"
	"path/filepath"
)

func readGameServerTemplates() []*gameServerTemplate {
	fmt.Println("Parsing for GameServerTemplates")
	templatesFolder := "./gameservertemplates"
	files, err := ioutil.ReadDir(templatesFolder)
	if err != nil {
		log.Fatal(err)
	}
	templates := make([]*gameServerTemplate, 0)
	for _, file := range files {
		if file.Mode().IsDir() {
			folder := file
			fmt.Println("Testing Template: " + folder.Name())
			tmpl, err := parseGameServerTemplate(templatesFolder, folder.Name())
			if err != nil {
				fmt.Println("Failed Processing: " + folder.Name() + "(" + err.Error() + ")")
			} else {
				templates = append(templates, tmpl)
				fmt.Println("Parsed: " + folder.Name())

			}
		}
	}
	return templates
}

func parseGameServerTemplate(basePath string, folder string) (*gameServerTemplate, error) {
	templatePath := filepath.Join(basePath, folder)
	//FIXME filepath in stdout is incoherent
	cfg, err := parseConfigMap(templatePath)
	if err != nil {
		return nil, err
	}
	pvc, err := parsePVC(templatePath)
	if err != nil {
		return nil, err
	}
	depl, err := parseDeployment(templatePath)
	if err != nil {
		return nil, err
	}
	svc, err := parseService(templatePath)
	if err != nil {
		return nil, err
	}
	gst := gameServerTemplate{folder, *cfg, *pvc, *depl, *svc}
	err = gst.Validate("GameServerTemplate")
	if err != nil {
		return nil, err
	}
	err = gst.Rename()
	if err != nil {
		return nil, err
	}
	return &gst, nil
}

func parseConfigMap(templatePath string) (*kubernetesComponentConfigMap, error) {
	filePath := templatePath + "/0-configmap.yaml"
	kubernetesTemplate, err := parseKubernetesTemplate(filePath, "ConfigMap")
	if parsed, ok := kubernetesTemplate.(*v1.ConfigMap); !ok {
		fmt.Println("Type Conversion Check failed!")
		return nil, err
	} else {
		fmt.Println("Parsed: " + filePath)
		return &kubernetesComponentConfigMap{*parsed}, nil
	}
}

func parsePVC(templatePath string) (*kubernetesComponentPVC, error) {
	filePath := templatePath + "/1-pvc.yaml"
	kubernetesTemplate, err := parseKubernetesTemplate(filePath, "PersistentVolumeClaim")
	if parsed, ok := kubernetesTemplate.(*v1.PersistentVolumeClaim); !ok {
		fmt.Println("Type Conversion Check failed!")
		return nil, err
	} else {
		fmt.Println("Parsed: " + filePath)
		return &kubernetesComponentPVC{*parsed}, nil
	}
}

func parseDeployment(templatePath string) (*kubernetesComponentDeployment, error) {
	filePath := templatePath + "/2-deployment.yaml"
	kubernetesTemplate, err := parseKubernetesTemplate(filePath, "Deployment.apps")
	if parsed, ok := kubernetesTemplate.(*appsv1.Deployment); !ok {
		fmt.Println("Type Conversion Check failed!")
		return nil, err
	} else {
		fmt.Println("Parsed: " + filePath)
		return &kubernetesComponentDeployment{*parsed}, nil
	}
}

func parseService(templatePath string) (*kubernetesComponentService, error) {
	filePath := templatePath + "/3-service.yaml"
	kubernetesTemplate, err := parseKubernetesTemplate(filePath, "Service")
	if parsed, ok := kubernetesTemplate.(*v1.Service); !ok {
		fmt.Println("Type Conversion Check failed!")
		return nil, err
	} else {
		fmt.Println("Parsed: " + filePath)
		return &kubernetesComponentService{*parsed}, nil
	}
}

func parseKubernetesTemplate(filePath string, kind string) (runtime.Object, error) {
	parseFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Could not read: " + filePath)
		return nil, err
	}
	kubernetesTemplate, parsedGroupVersionKind, err := sch.Codecs.UniversalDeserializer().Decode(parseFile, nil, nil)
	if err != nil {
		fmt.Println("Parsing failed")
		return nil, err
	}
	_, groupkind := schema.ParseKindArg(kind)
	if !cmp.Equal(groupkind.WithVersion("v1"), *parsedGroupVersionKind) {
		errMsg := "Invalid Type in File! Expecting: " + kind + " Found: " + parsedGroupVersionKind.Kind
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	return kubernetesTemplate, nil
}
