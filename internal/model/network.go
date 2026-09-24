package model

import "time"

type NetworkSpec struct {
	Region string
}

func NewNetworkSpec(region string) NetworkSpec {
	return NetworkSpec{}
}

func (s NetworkSpec) Validate() error {
	return nil
}

type Network struct {
	NetworkName  string
	CIDR         string
	SpecData     NetworkSpec
	ResourceType string
	DependsOn    []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewNetwork(networkName string, networkSpec NetworkSpec) Network {
	return Network{
		NetworkName:  networkName,
		ResourceType: ResourceTypeNetwork,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (n *Network) Name() string { return n.NetworkName }

func (n *Network) Key() string {
	return "network/" + n.NetworkName
}
func (n *Network) Dependencies() []string {
	return []string{}
}
func (n *Network) Type() string {
	return ResourceTypeNetwork
}
func (n *Network) Spec() Spec {
	return n.SpecData
}

func (n *Network) SetUpdatedAt(time.Time) {

}
