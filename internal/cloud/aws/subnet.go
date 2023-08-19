package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a *CloudProvider) retrieveSubnet(subnetName string) string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(subnetName)},
			},
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error describing Subnets: ", err)
	}

	if len(output.Subnets) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", subnetName, *output.Subnets[0].SubnetId)
	return *output.Subnets[0].SubnetId
}

func (a *CloudProvider) createSubnet(vpcId *string, subnetName string, subnetCidr string, subnetAz string) string {
	svc := ec2.New(a.session)

	subnet, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            vpcId,
		CidrBlock:        aws.String(subnetCidr),
		AvailabilityZone: aws.String(subnetAz),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("subnet"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(subnetName),
				},
			},
		}},
	})

	if err != nil {
		log.Fatal("[🐶] Error creating Subnet: ", err)
	}

	fmt.Println("[🐶] Subnet creation requested, waiting for completion...")
	if err := svc.WaitUntilSubnetAvailable(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(subnetName)},
			},
		},
	}); err != nil {
		log.Fatal("[🐶] Error waiting for Subnet creation: ", err)
	}
	fmt.Printf("[🐶] %s Successfully created: %s\n", subnetName, *subnet.Subnet.SubnetId)

	return *subnet.Subnet.SubnetId
}

func (a *CloudProvider) addPublicIpAutoAssignToSubnet(subnetId string) string {
	svc := ec2.New(a.session)

	_, err := svc.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
		SubnetId: aws.String(subnetId),
		MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error modifying Subnet: ", err)
	}

	fmt.Printf("[🐶] Successfully applied Public Ip Auto Assign to subnet %s\n", a.awsConfig.VPC.Subnets.PublicSubnetName)
	return ""
}
