package cloud

import (
	"github.com/dumunari/spikectl/internal/cloud/aws"
	"github.com/dumunari/spikectl/internal/cloud/azure"
	gcp "github.com/dumunari/spikectl/internal/cloud/gcp"
	"github.com/dumunari/spikectl/internal/config"
)

type CloudProvider interface {
	InstantiateKubernetesCluster() error
}

func NewCloudProvider(cfg *config.SpikeConfig) CloudProvider {
	if cfg.IDP.CloudProvider == config.AZURE {
		return azure.NewAzureCloudProvider(cfg)
	}
	if cfg.IDP.CloudProvider == config.GCP {
		return gcp.NewGcpCloudProvider(cfg)
	}
	return aws.NewAwsCloudProvider(cfg)
}
