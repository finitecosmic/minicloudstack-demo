package model

import "net/http"

type ObjectHandler interface {
	Put(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type ObjectStoreConfig struct {
	Name       string
	Region     string
	Versioning bool
	Encryption string
	Tags       map[string]string
}

func NewObjectStore(c ObjectStoreConfig) *ObjectStoreConfig {
	return &ObjectStoreConfig{}
}

type ObjectStore struct {
	ObjectStoreConfig
}

func (b *ObjectStoreConfig) Error() string {
	//TODO implement me
	panic("implement me")
}
