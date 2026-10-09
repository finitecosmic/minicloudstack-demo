package model

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
