package gcp

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

func (a *CloudProvider) getOrCreateVpc(vpcName string) string {
	vpcId := a.retrieveVpc()

	if vpcId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", vpcName)
		vpcId = a.createVpc()
	}
	return vpcId
}

func (a *CloudProvider) retrieveVpc() string {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal("[🐶] Error creating Compute Engine Service: ", err)
	}

	filter := fmt.Sprintf("name eq %s", a.gcpConfig.VPC.Name)
	resp, err := service.Networks.List(a.gcpConfig.ProjectId).Filter(filter).Context(ctx).Do()
	if err != nil {
		log.Fatal("[🐶] Error listing VPCs: ", err)
	}

	if len(resp.Items) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %d\n", a.gcpConfig.VPC.Name, resp.Items[0].Id)
	return resp.Items[0].SelfLink
}

func (a *CloudProvider) createVpc() string {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal("[🐶] Error creating Compute Engine Service: ", err)
	}

	vpc := &compute.Network{
		Name:                  a.gcpConfig.VPC.Name,
		AutoCreateSubnetworks: false,
		ForceSendFields:       []string{"AutoCreateSubnetworks"},
	}

	op, err := service.Networks.Insert(a.gcpConfig.ProjectId, vpc).Context(ctx).Do()
	if err != nil {
		log.Fatal("[🐶] Error creating VPC: ", err)
	}

	wait_op, err := service.GlobalOperations.Wait(a.gcpConfig.ProjectId, op.Name).Context(ctx).Do()

	if err != nil || wait_op.Error != nil {
		log.Fatal("[🐶] Error waiting for operation: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", a.gcpConfig.VPC.Name, op.TargetLink)

	return op.TargetLink
}
