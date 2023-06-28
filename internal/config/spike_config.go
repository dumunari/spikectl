package config

type ResourceGroupConfig struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type AzureConfig struct {
	SubscriptionId      string              `json:"subscription_id"`
	ResourceGroupConfig ResourceGroupConfig `json:"resource_group"`
}

type IDPConfig struct {
	CloudProvider CloudProvider `json:"cloud_provider"`
	AzureConfig   AzureConfig   `json:"azure"`
}

type SpikeConfig struct {
	IDP IDPConfig `json:"idp"`
}
