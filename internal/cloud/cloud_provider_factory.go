package cloud

import "spikectl/internal/config"

type CloudProvider interface {
	CreateResourceGroup()
	CreateKubernetesCluster()
}

func NewCloudProvider(cfg *config.SpikeConfig) CloudProvider {
	if cfg.IDP.CloudProvider == config.AZURE {
		return NewAzureCloudProvider(cfg)
	}
	return NewAwsCloudProvider()
}
