package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

func (a *CloudProvider) retrieveNodeGroup() string {
	svc := eks.New(a.session)

	output, err := svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(a.awsConfig.EKS.Name),
		NodegroupName: aws.String(a.awsConfig.EKS.NodeGroup.Name),
	})

	if err != nil {
		fmt.Printf("[🐶] Error describing NodeGroup: %s\n", err)
		return ""
	}

	fmt.Printf("[🐶] Found %s with Arn: %s\n", a.awsConfig.EKS.NodeGroup.Name, *output.Nodegroup.NodegroupArn)
	return *output.Nodegroup.NodegroupArn
}

func (a *CloudProvider) createNodeGroup(nodeRole string, subnetIds ...string) {
	svc := eks.New(a.session)

	_, err := svc.CreateNodegroup(&eks.CreateNodegroupInput{
		ClusterName:   aws.String(a.awsConfig.EKS.Name),
		NodegroupName: aws.String(a.awsConfig.EKS.NodeGroup.Name),
		Subnets:       aws.StringSlice(subnetIds),
		InstanceTypes: aws.StringSlice(a.awsConfig.EKS.NodeGroup.InstanceType),
		NodeRole:      aws.String(nodeRole),
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(a.awsConfig.EKS.NodeGroup.Scaling.DesiredSize),
			MaxSize:     aws.Int64(a.awsConfig.EKS.NodeGroup.Scaling.MaxSize),
			MinSize:     aws.Int64(a.awsConfig.EKS.NodeGroup.Scaling.MinSize),
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error creating NodeGroup: ", err)
	}

	//TODO: too messy
	fmt.Println("[🐶] NodeGroup creation requested, waiting for completion...")
	if err := svc.WaitUntilNodegroupActive(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(a.awsConfig.EKS.Name),
		NodegroupName: aws.String(a.awsConfig.EKS.NodeGroup.Name),
	}); err != nil {
		log.Fatal("[🐶] Error waiting for NodeGroup creation: ", err)
	}
	fmt.Printf("[🐶] Successfully created %s\n", a.awsConfig.EKS.NodeGroup.Name)
}
