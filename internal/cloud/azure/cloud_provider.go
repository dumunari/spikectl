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

<<<<<<< HEAD
	//_, err = p.retrieveVirtualNetwork(rg)

	_, err = p.createOrUpdateAKS(&AksParameters{
		ResourceGroup:      rg,
		ManagedClusterName: "cluster-test",
	})

	if err != nil {
		return err
	}

	return nil
=======
	az.retrieveVirtualNetwork(rg)

	return config.KubeConfig{}
>>>>>>> bf6e254 (feat: bind core components installation to aks and eks clusters)
}
