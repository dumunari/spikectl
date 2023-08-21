package gcp

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

func (g *CloudProvider) getOrCreateVpc(vpcName string) string {
	vpcId := g.retrieveVpc()

	if vpcId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", vpcName)
		vpcId = g.createVpc()
	}
	return vpcId
}

func (g *CloudProvider) retrieveVpc() string {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal("[🐶] Error creating Compute Engine Service: ", err)
	}

	filter := fmt.Sprintf("name eq %s", g.gcpConfig.VPC.Name)
	resp, err := service.Networks.List(g.gcpConfig.ProjectId).Filter(filter).Context(ctx).Do()
	if err != nil {
		log.Fatal("[🐶] Error listing VPCs: ", err)
	}

	if len(resp.Items) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %d\n", g.gcpConfig.VPC.Name, resp.Items[0].Id)
	return resp.Items[0].SelfLink
}

func (g *CloudProvider) createVpc() string {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal("[🐶] Error creating Compute Engine Service: ", err)
	}

	vpc := &compute.Network{
		Name:                  g.gcpConfig.VPC.Name,
		AutoCreateSubnetworks: false,
		ForceSendFields:       []string{"AutoCreateSubnetworks"},
	}

	op, err := service.Networks.Insert(g.gcpConfig.ProjectId, vpc).Context(ctx).Do()
	if err != nil {
		log.Fatal("[🐶] Error creating VPC: ", err)
	}

	wait_op, err := service.GlobalOperations.Wait(g.gcpConfig.ProjectId, op.Name).Context(ctx).Do()

	if err != nil || wait_op.Error != nil {
		log.Fatal("[🐶] Error waiting for operation: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", g.gcpConfig.VPC.Name, op.TargetLink)

	return op.TargetLink
}
