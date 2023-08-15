package cloud

import (
	"github.com/dumunari/spikectl/internal/cloud/aws"
	"github.com/dumunari/spikectl/internal/cloud/azure"
	"github.com/dumunari/spikectl/internal/cloud/gcp"
	"github.com/dumunari/spikectl/internal/config"
)

type CloudProvider interface {
	InstantiateKubernetesCluster() error
}

func NewCloudProvider(cfg *config.Spike) CloudProvider {
	if cfg.Spike.CloudProvider == config.AZURE {
		return azure.NewAzureCloudProvider(cfg)
	}
	if cfg.Spike.CloudProvider == config.GCP {
		return gcp.NewGcpCloudProvider(cfg)
	}
	return aws.NewAwsCloudProvider(cfg)
}
