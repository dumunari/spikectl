package cloud

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"log"
	"spikectl/internal/config"
)

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type AzureCloudProvider struct {
	azureConfig config.AzureConfig
}

func NewAzureCloudProvider(config *config.SpikeConfig) *AzureCloudProvider {
	cloudProvider := AzureCloudProvider{}
	cloudProvider.azureConfig = config.IDP.AzureConfig

	return &cloudProvider
}

func (p *AzureCloudProvider) CreateKubernetesCluster() {}

func (p *AzureCloudProvider) CreateResourceGroup() {
	subscriptionID := p.azureConfig.SubscriptionId
	if len(subscriptionID) == 0 {
		log.Fatal("Azure Subscription ID wasn't provided")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	resourcesClientsFactory, err := armresources.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	resourceGroupClient := resourcesClientsFactory.NewResourceGroupsClient()

	exists, err := checkExistenceResourceGroup(resourceGroupClient, ctx, p.azureConfig.ResourceGroupConfig.Name)

	if !exists {
		resourceGroup, err := createResourceGroup(resourceGroupClient, ctx, p.azureConfig.ResourceGroupConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("ResourceGroup created. ID: %s", *resourceGroup.ID)
	}
}

func createResourceGroup(client *armresources.ResourceGroupsClient, ctx context.Context, cfg config.ResourceGroupConfig) (*armresources.ResourceGroup, error) {
	resourceGroupResp, err := client.CreateOrUpdate(
		ctx,
		cfg.Name,
		armresources.ResourceGroup{
			Location: to.Ptr(cfg.Location),
		},
		nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func checkExistenceResourceGroup(client *armresources.ResourceGroupsClient, ctx context.Context, name string) (bool, error) {
	boolResp, err := client.CheckExistence(ctx, name, nil)
	if err != nil {
		return false, err
	}

	return boolResp.Success, nil
}
