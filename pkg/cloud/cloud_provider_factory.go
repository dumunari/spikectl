package cloud

import "spikectl/pkg/config"

type CloudProvider interface {
	CreateKubernetesCluster()
}

type AzureCloudProvider struct {
}

func NewAzureCloudProvider() *AzureCloudProvider {
	return &AzureCloudProvider{}
}

func (p *AzureCloudProvider) CreateKubernetesCluster() {
}

type AwsCloudProvider struct {
}

func CreateCloudProvider(cfg config.SpikeConfig) CloudProvider {
	if cfg.IDP.CloudProvider == config.AZURE {
		return NewAzureCloudProvider()
	}
}
