package azure

import (
	"log"
	"spikectl/internal/config"
)

type AzureCloudProvider struct {
	azureConfig config.AzureConfig
}

func NewAzureCloudProvider(config *config.SpikeConfig) *AzureCloudProvider {
	cloudProvider := AzureCloudProvider{}
	cloudProvider.azureConfig = config.IDP.AzureConfig

	if len(cloudProvider.azureConfig.SubscriptionId) == 0 {
		log.Fatal("Subscription id wasn't provided")
	}

	return &cloudProvider
}

func (p *AzureCloudProvider) InstantiateKubernetesCluster() error {
	rg, err := p.retrieveResourceGroup()
	if err != nil {
		return err
	}

	vnet, err := p.retrieveVNet(rg)

	return nil
}
