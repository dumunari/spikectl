package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a *CloudProvider) retrieveVpc() string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(a.awsConfig.VPC.Name)},
			},
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error describing VPCs: ", err)
	}

	if len(output.Vpcs) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", a.awsConfig.VPC.Name, *output.Vpcs[0].VpcId)

	return *output.Vpcs[0].VpcId
}

func (a *CloudProvider) createVpc() string {
	svc := ec2.New(a.session)

	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(a.awsConfig.VPC.CIDR),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("vpc"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.Name),
				},
			},
		}},
	}

	vpc, err := svc.CreateVpc(input)

	if err != nil {
		log.Fatal("[🐶] Error creating VPC: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", a.awsConfig.VPC.Name, *vpc.Vpc.VpcId)

	return *vpc.Vpc.VpcId
}

func (a *CloudProvider) retrieveMainRouteTable(vpcId string) string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error describing Route Tables: ", err)
	}

	if len(output.RouteTables) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found Route Table with Id: %s\n", *output.RouteTables[0].RouteTableId)

	return *output.RouteTables[0].RouteTableId
}

func (a *CloudProvider) createPublicRouteTable(vpcId string) string {
	svc := ec2.New(a.session)

	output, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("route-table"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.PublicRouteTable.Name),
				},
			},
		}},
	})

	if err != nil {
		log.Fatal("[🐶] Error creating Route Table: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created with Id: %s\n", a.awsConfig.VPC.PublicRouteTable.Name, *output.RouteTable.RouteTableId)

	return *output.RouteTable.RouteTableId
}

func (a *CloudProvider) retrievePublicRouteTable(vpcId string) string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(a.awsConfig.VPC.PublicRouteTable.Name)},
			},
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error describing Route Tables: ", err)
	}

	if len(output.RouteTables) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", a.awsConfig.VPC.PublicRouteTable.Name, *output.RouteTables[0].RouteTableId)

	return *output.RouteTables[0].RouteTableId
}
