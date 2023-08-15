package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a CloudProvider) retrieveNatGateway() string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(a.awsConfig.VPC.NatGateway.Name)},
			},
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error describing NAT Gateway: ", err)
	}

	if len(output.NatGateways) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", a.awsConfig.VPC.NatGateway.Name, *output.NatGateways[0].NatGatewayId)
	return *output.NatGateways[0].NatGatewayId
}

func (a CloudProvider) createNatGateway(publicSubnetId string) string {
	svc := ec2.New(a.session)

	eipID := a.createIPAllocation()

	input := &ec2.CreateNatGatewayInput{
		AllocationId: aws.String(eipID),
		SubnetId:     aws.String(publicSubnetId),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("natgateway"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.NatGateway.Name),
				},
			},
		}},
	}

	igw, err := svc.CreateNatGateway(input)

	if err != nil {
		log.Fatal("[🐶] Error creating NAT Gateway: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", a.awsConfig.VPC.NatGateway.Name, *igw.NatGateway.NatGatewayId)
	return *igw.NatGateway.NatGatewayId
}

func (a CloudProvider) createIPAllocation() string {
	svc := ec2.New(a.session)

	input := &ec2.AllocateAddressInput{
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("elastic-ip"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.NatGateway.Name),
				},
			},
		}},
	}

	eip, err := svc.AllocateAddress(input)

	if err != nil {
		log.Fatal("[🐶] Error allocating IP Address: ", err)
	}

	fmt.Printf("[🐶] Successfully allocated IP: %s with Allocation ID: %s\n", *eip.PublicIp, *eip.AllocationId)
	return *eip.AllocationId
}

func (a CloudProvider) addNatGatewayToVpcMainRouteTable(routeTableId string, ngwId string) string {
	svc := ec2.New(a.session)

	input := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(routeTableId),
		NatGatewayId:         aws.String(ngwId),
	}

	_, err := svc.CreateRoute(input)

	if err != nil {
		log.Fatal("[🐶] Error creating route: ", err)
	}

	fmt.Printf("[🐶] Successfully created Route\n")
	return ""
}
