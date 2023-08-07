package config

type AzureConfig struct {
	SubscriptionId       string               `json:"subscription_id"`
	ResourceGroupConfig  ResourceGroupConfig  `json:"resource_group"`
	VirtualNetworkConfig VirtualNetworkConfig `json:"virtual_network"`
}

type ResourceGroupConfig struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type Subnets struct {
	Name          string `json:"name"`
	AddressPrefix string `json:"address_prefix"`
}

type VirtualNetworkConfig struct {
	Name            string    `json:"name"`
	AddressPrefixes []string  `json:"address_prefixes"`
	Subnets         []Subnets `json:"subnets"`
}
