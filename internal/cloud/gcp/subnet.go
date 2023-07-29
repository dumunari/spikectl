package gcp

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

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

func (a *CloudProvider) createSubnet(vpcLink, subnetName, subnetCidr, subnetRegion string) string {
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

	op, err := service.Subnetworks.Insert(a.gcpConfig.ProjectId, subnetRegion, subnet).Context(ctx).Do()

	if err != nil {
		log.Fatal("[🐶] Error creating Subnet: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", subnetName, op.TargetLink)
	return op.TargetLink
}
