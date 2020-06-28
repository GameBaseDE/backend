package openapi

type kubernetesComponentWrapper interface {
	GetName() string
	Rename(newName string) error
	Validate(templatePrefix string) error
}
