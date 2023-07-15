package cloud

import (
	"github.com/dumunari/spikectl/internal/cloud/aws"
	"github.com/dumunari/spikectl/internal/cloud/azure"
	"github.com/dumunari/spikectl/internal/config"
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
