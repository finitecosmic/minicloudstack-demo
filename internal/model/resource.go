package model

type Plan struct {
	resources []Resource
}

type Resource interface {
	Key() string
	Name() string
	Dependencies() []string
	Type() string
	Spec() Spec
}
