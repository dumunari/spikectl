package config

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

type AgentPools struct {
	Name     string `json:"name"`
	VmsSize  string `json:"vms_size"`
	MaxPods  int    `json:"max_pods"`
	MinCount int    `json:"min_count"`
	Count    int    `json:"count"`
	MaxCount int    `json:"max_count"`
}

type ServiceProvider struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Aks struct {
	ClusterName     string          `json:"cluster_name"`
	AgentPools      []AgentPools    `json:"agent_pools"`
	ServiceProvider ServiceProvider `json:"service_provider"`
}

type AzureConfig struct {
	SubscriptionId       string               `json:"subscription_id"`
	ResourceGroupConfig  ResourceGroupConfig  `json:"resource_group"`
	VirtualNetworkConfig VirtualNetworkConfig `json:"virtual_network"`
	Aks                  Aks                  `json:"aks"`
}
