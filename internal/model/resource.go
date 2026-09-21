package model

type Plan struct {
	resources []Resource
}

type Spec interface{}

type Resource interface {
	Key() string
	Name() string
	GetSpec() Spec
	GetDependencies() []string
	ResourceType() string
}
