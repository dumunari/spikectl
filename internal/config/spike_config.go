package config

type SpikeConfig struct {
	IDP IDPConfig `json:"idp"`
}

type IDPConfig struct {
	CloudProvider CloudProvider `json:"cloud_provider"`
	AzureConfig   AzureConfig   `json:"azure"`
	AwsConfig     AwsConfig     `json:"aws"`
}

type AwsConfig struct {
	Profile string    `json:"profile"`
	Region  string    `json:"region"`
	VPC     VPCConfig `json:"vpc"`
}

type VPCConfig struct {
	Name    string       `json:"name"`
	CIDR    string       `json:"cidr"`
	Subnets SubnetConfig `json:"subnets"`
}

type SubnetConfig struct {
	PublicSubnetName  string `json:"public_subnet_name"`
	PublicSubnetCidr  string `json:"public_subnet_cidr"`
	PrivateSubnetName string `json:"private_subnet_name"`
	PrivateSubnetCidr string `json:"private_subnet_cidr"`
}

type AzureConfig struct {
	SubscriptionId       string               `json:"subscription_id"`
	ResourceGroupConfig  ResourceGroupConfig  `json:"resource_group"`
	VirtualNetworkConfig VirtualNetworkConfig `json:"virtual_network"`
}

type ResourceGroupConfig struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type VirtualNetworkConfig struct {
	Name            string    `json:"name"`
	AddressPrefixes []string  `json:"address_prefixes"`
	Subnets         []Subnets `json:"subnets"`
}

type Subnets struct {
	Name          string `json:"name"`
	AddressPrefix string `json:"address_prefix"`
}
