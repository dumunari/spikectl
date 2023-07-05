package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"log"
	"spikectl/internal/config"
)

type CloudProvider struct {
	azureConfig config.AzureConfig
	credentials *azidentity.DefaultAzureCredential
}

func NewAzureCloudProvider(config *config.SpikeConfig) *CloudProvider {
	cloudProvider := CloudProvider{}
	cloudProvider.azureConfig = config.IDP.AzureConfig

	if len(cloudProvider.azureConfig.SubscriptionId) == 0 {
		log.Fatal("Subscription id wasn't provided")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	cloudProvider.credentials = cred
	return &cloudProvider
}

func (p *CloudProvider) InstantiateKubernetesCluster() error {
	rg, err := p.retrieveResourceGroup()
	if err != nil {
		return err
	}

	_, err = p.retrieveVirtualNetwork(rg)
	return nil
}
