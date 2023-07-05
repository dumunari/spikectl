package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"log"
	"spikectl/internal/config"
)

func (p *CloudProvider) retrieveResourceGroup() (*armresources.ResourceGroup, error) {

	subscriptionID := p.azureConfig.SubscriptionId
	if len(subscriptionID) == 0 {
		log.Fatal("Azure Subscription ID wasn't provided")
	}

	ctx := context.Background()
	resourcesClientsFactory, err := armresources.NewClientFactory(subscriptionID, p.credentials, nil)
	if err != nil {
		log.Fatal(err)
	}

	resourceGroupClient := resourcesClientsFactory.NewResourceGroupsClient()

	exists, err := checkExistenceResourceGroup(resourceGroupClient, ctx, p.azureConfig.ResourceGroupConfig)

	var resourceGroup *armresources.ResourceGroup

	if !exists {
		resourceGroup, err = createResourceGroup(resourceGroupClient, ctx, p.azureConfig.ResourceGroupConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("ResourceGroup created. ID: %s", *resourceGroup.ID)
	} else {
		resourceGroup, err = getResourceGroup(resourceGroupClient, ctx, p.azureConfig.ResourceGroupConfig)
	}

	return resourceGroup, nil

}

func getResourceGroup(client *armresources.ResourceGroupsClient, ctx context.Context, groupConfig config.ResourceGroupConfig) (*armresources.ResourceGroup, error) {
	resourceGroupResp, err := client.Get(ctx, groupConfig.Name, nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
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

func checkExistenceResourceGroup(client *armresources.ResourceGroupsClient, ctx context.Context, cfg config.ResourceGroupConfig) (bool, error) {
	boolResp, err := client.CheckExistence(ctx, cfg.Name, nil)
	if err != nil {
		return false, err
	}

	return boolResp.Success, nil
}
