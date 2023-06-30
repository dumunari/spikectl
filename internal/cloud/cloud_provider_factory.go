package cloud

import (
	"spikectl/internal/cloud/aws"
	"spikectl/internal/cloud/azure"
	"spikectl/internal/config"
)

type CloudProvider interface {
	InstantiateKubernetesCluster() error
}

func NewCloudProvider(cfg *config.SpikeConfig) CloudProvider {
	if cfg.IDP.CloudProvider == config.AZURE {
		return azure.NewAzureCloudProvider(cfg)
	}
	return aws.NewAwsCloudProvider()
}
