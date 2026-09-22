package model

type Spec interface {
	Validate() error
}
