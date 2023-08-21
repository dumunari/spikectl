package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a *CloudProvider) retrieveNatGateway() string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(a.awsConfig.VPC.NatGateway.Name)},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
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

func (a *CloudProvider) createNatGateway(publicSubnetId string) string {
	svc := ec2.New(a.session)

	eipID := a.createIPAllocation()

	igw, err := svc.CreateNatGateway(&ec2.CreateNatGatewayInput{
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
	})

	if err != nil {
		log.Fatal("[🐶] Error creating NAT Gateway: ", err)
	}

	fmt.Println("[🐶] NAT Gateway creation requested, waiting for completion...")
	if err := svc.WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(a.awsConfig.VPC.NatGateway.Name)},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
		},
	}); err != nil {
		log.Fatal("[🐶] Error waiting for NAT Gateway creation: ", err)
	}
	fmt.Printf("[🐶] %s Successfully created: %s\n", a.awsConfig.VPC.NatGateway.Name, *igw.NatGateway.NatGatewayId)

	return *igw.NatGateway.NatGatewayId
}

func (a *CloudProvider) createIPAllocation() string {
	svc := ec2.New(a.session)

	eip, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("elastic-ip"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.NatGateway.Name),
				},
			},
		}},
	})

	if err != nil {
		log.Fatal("[🐶] Error allocating IP Address: ", err)
	}

	fmt.Printf("[🐶] Successfully allocated IP: %s with Allocation ID: %s\n", *eip.PublicIp, *eip.AllocationId)
	return *eip.AllocationId
}

func (a *CloudProvider) addNatGatewayToVpcMainRouteTable(routeTableId string, ngwId string) string {
	svc := ec2.New(a.session)

	_, err := svc.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(routeTableId),
		NatGatewayId:         aws.String(ngwId),
	})

	if err != nil {
		log.Fatal("[🐶] Error creating route: ", err)
	}

	fmt.Printf("[🐶] Successfully created Route for %s Main Route Table <-> %s association\n", a.awsConfig.VPC.Name, a.awsConfig.VPC.NatGateway.Name)
	return ""
}
