package azure

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/dumunari/spikectl/internal/config"
)

type CloudProvider struct {
	azureConfig config.AzureConfig
	credentials *azidentity.DefaultAzureCredential
}

func NewAzureCloudProvider(config *config.Spike) *CloudProvider {
	cloudProvider := CloudProvider{}
	cloudProvider.azureConfig = config.Spike.AzureConfig

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

func (az *CloudProvider) InstantiateKubernetesCluster() config.KubeConfig {
	rg, _ := az.retrieveResourceGroup()

	_, err := az.retrieveVirtualNetwork(rg)

	aksCluster, err := az.createOrUpdateAKS(&AksParameters{
		ResourceGroup:      rg,
		ManagedClusterName: "cluster-test",
	})

	if err != nil {
		log.Fatal("Error creating cluster:", err)
	}

	return az.retrieveKubeConfigInfo(*aksCluster)
}
