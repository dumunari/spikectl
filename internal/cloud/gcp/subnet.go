package gcp

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

func (a *CloudProvider) getOrCreateSubnet(vpcLink, subnetName, subnetRegion, subnetCidr string) string {
	subnetId := a.retrieveSubnet(subnetName, subnetRegion)

	if subnetId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", subnetName)
		subnetId = a.createSubnet(vpcLink, subnetName, subnetCidr, subnetRegion)
	}

	return subnetId
}

func (a *CloudProvider) retrieveSubnet(subnetName, subnetRegion string) string {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal("[🐶] Error creating Compute Engine Service: ", err)
	}

	filter := fmt.Sprintf("name eq %s", subnetName)
	resp, err := service.Subnetworks.List(a.gcpConfig.ProjectId, subnetRegion).Filter(filter).Context(ctx).Do()

	if err != nil {
		log.Fatal("[🐶] Error listing Subnets: ", err)
	}

	if len(resp.Items) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %d\n", subnetName, resp.Items[0].Id)
	return resp.Items[0].SelfLink
}

func (g *CloudProvider) createSubnet(vpcLink, subnetName, subnetCidr, subnetRegion string) string {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal("[🐶] Error creating Compute Engine Service: ", err)
	}

	subnet := &compute.Subnetwork{
		Name:        subnetName,
		IpCidrRange: subnetCidr,
		Network:     vpcLink,
		Region:      subnetRegion,
	}

	op, err := service.Subnetworks.Insert(g.gcpConfig.ProjectId, subnetRegion, subnet).Context(ctx).Do()

	if err != nil {
		log.Fatal("[🐶] Error creating Subnet: ", err)
	}

	wait_op, err := service.RegionOperations.Wait(g.gcpConfig.ProjectId, subnetRegion, op.Name).Context(ctx).Do()

	if err != nil || wait_op.Error != nil {
		log.Fatal("[🐶] Error waiting for operation: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", subnetName, op.TargetLink)
	return op.TargetLink
}
