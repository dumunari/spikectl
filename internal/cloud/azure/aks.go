package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"log"
)

type AksParameters struct {
	ResourceGroup      *armresources.ResourceGroup
	ManagedClusterName string
}

func (p *CloudProvider) createOrUpdateAKS(params *AksParameters) (*armcontainerservice.ManagedCluster, error) {

	subscriptionID := p.azureConfig.SubscriptionId

	containerServiceClientFactory, err := armcontainerservice.NewClientFactory(subscriptionID, p.credentials, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	managedClusterClient := containerServiceClientFactory.NewManagedClustersClient()

	pollerResp, err := managedClusterClient.BeginCreateOrUpdate(
		ctx,
		*params.ResourceGroup.Name,
		params.ManagedClusterName,
		armcontainerservice.ManagedCluster{
			Location: params.ResourceGroup.Location,
			Properties: &armcontainerservice.ManagedClusterProperties{
				AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
					{
						Name:              to.Ptr("default"),
						Count:             to.Ptr[int32](1),
						VMSize:            to.Ptr("Standard_DS2_v2"),
						MaxPods:           to.Ptr[int32](110),
						MinCount:          to.Ptr[int32](1),
						MaxCount:          to.Ptr[int32](10),
						OSType:            to.Ptr(armcontainerservice.OSTypeLinux),
						Type:              to.Ptr(armcontainerservice.AgentPoolTypeVirtualMachineScaleSets),
						EnableAutoScaling: to.Ptr(true),
						Mode:              to.Ptr(armcontainerservice.AgentPoolModeSystem),
					},
				},
				DNSPrefix: to.Ptr("idp"),
				ServicePrincipalProfile: &armcontainerservice.ManagedClusterServicePrincipalProfile{
					ClientID: &p.azureConfig.Aks.ServiceProvider.ClientID,
					Secret:   &p.azureConfig.Aks.ServiceProvider.ClientSecret,
				},
			},
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedCluster, nil
}
