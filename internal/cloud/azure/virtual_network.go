package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"log"
	"spikectl/internal/config"
)

func (p *CloudProvider) retrieveVirtualNetwork(rg *armresources.ResourceGroup) (*armnetwork.VirtualNetwork, error) {
	networkClientFactory, err := armnetwork.NewClientFactory(p.azureConfig.SubscriptionId, p.credentials, nil)
	if err != nil {
		return nil, err
	}

	client := networkClientFactory.NewVirtualNetworksClient()

	ctx := context.Background()
	if vnet, ok := checkExistenceVirtualNetwork(ctx, client, rg, p.azureConfig.VirtualNetworkConfig); ok {
		return vnet, nil
	}

	vnet, err := createVirtualNetwork(ctx, client, rg, p.azureConfig.VirtualNetworkConfig)
	if err != nil {
		return nil, err
	}

	return vnet, nil
}

func checkExistenceVirtualNetwork(ctx context.Context, client *armnetwork.VirtualNetworksClient, resourceGroup *armresources.ResourceGroup, vnetConfig config.VirtualNetworkConfig) (*armnetwork.VirtualNetwork, bool) {

	resp, err := client.Get(ctx, *resourceGroup.Name, vnetConfig.Name, nil)
	if err != nil {
		log.Println(fmt.Errorf("error trying to retrieve the virtual network: %v", err))
		return nil, false
	}

	return &resp.VirtualNetwork, true
}

func createVirtualNetwork(ctx context.Context, client *armnetwork.VirtualNetworksClient, resourceGroup *armresources.ResourceGroup, vnetConfig config.VirtualNetworkConfig) (*armnetwork.VirtualNetwork, error) {

	addressPrefixes := createVNetAddressPrefixStructure(vnetConfig)

	subnets := createSubnetStructure(vnetConfig)

	vnet := armnetwork.VirtualNetwork{
		Location: resourceGroup.Location,
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: addressPrefixes,
			},
			Subnets: subnets,
		},
	}

	pollerResp, err := client.BeginCreateOrUpdate(
		ctx,
		*resourceGroup.Name,
		vnetConfig.Name,
		vnet,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualNetwork, nil
}

func createVNetAddressPrefixStructure(vnetConfig config.VirtualNetworkConfig) []*string {
	addressPrefixes := make([]*string, len(vnetConfig.AddressPrefixes))
	for i := range vnetConfig.AddressPrefixes {
		addressPrefixes[i] = &vnetConfig.AddressPrefixes[i]
	}
	return addressPrefixes
}

func createSubnetStructure(vnetConfig config.VirtualNetworkConfig) []*armnetwork.Subnet {
	subnets := make([]*armnetwork.Subnet, len(vnetConfig.Subnets))

	for i := range vnetConfig.Subnets {
		subnets[i] = &armnetwork.Subnet{
			Name: to.Ptr(vnetConfig.Subnets[i].Name),
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.Ptr(vnetConfig.Subnets[i].AddressPrefix),
			},
		}
	}
	return subnets
}
