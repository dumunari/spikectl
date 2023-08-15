package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a CloudProvider) retrieveInternetGateway() string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(a.awsConfig.VPC.InternetGateway.Name)},
			},
		},
	})

	if err != nil {
		log.Fatal("[🐶] Error describing Internet Gateway: ", err)
	}

	if len(output.InternetGateways) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", a.awsConfig.VPC.InternetGateway.Name, *output.InternetGateways[0].InternetGatewayId)
	return *output.InternetGateways[0].InternetGatewayId
}

func (a CloudProvider) createInternetGateway() string {
	svc := ec2.New(a.session)

	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("internet-gateway"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.InternetGateway.Name),
				},
			},
		}},
	}

	igw, err := svc.CreateInternetGateway(input)

	if err != nil {
		log.Fatal("[🐶] Error creating Internet Gateway: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", a.awsConfig.VPC.InternetGateway.Name, *igw.InternetGateway.InternetGatewayId)
	return *igw.InternetGateway.InternetGatewayId
}

func (a CloudProvider) attachInternetGateway(vpcId string, igwId string) string {
	svc := ec2.New(a.session)

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	}

	_, err := svc.AttachInternetGateway(input)

	if err != nil {
		log.Fatal("[🐶] Error attaching Internet Gateway to VPC: ", err)
	}

	fmt.Printf("[🐶] Successfully attached %s to %s\n", a.awsConfig.VPC.InternetGateway.Name, a.awsConfig.VPC.Name)
	return ""
}

func (a CloudProvider) addInternetGatewayToVpcPublicRouteTable(igwId string, publicRouteTableId string, publicSubnetId string) string {
	svc := ec2.New(a.session)

	input := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(publicRouteTableId),
		GatewayId:            aws.String(igwId),
	}

	_, err := svc.CreateRoute(input)

	if err != nil {
		log.Fatal("[🐶] Error creating route: ", err)
	}

	fmt.Printf("[🐶] Successfully created Route\n")
	return ""
}

func (a CloudProvider) associatePublicSubnetToPublicRouteTable(publicSubnetId string, publicRouteTableId string) string {
	svc := ec2.New(a.session)

	input := &ec2.AssociateRouteTableInput{
		SubnetId:     aws.String(publicSubnetId),
		RouteTableId: aws.String(publicRouteTableId),
	}

	_, err := svc.AssociateRouteTable(input)

	if err != nil {
		log.Fatal("[🐶] Error associating Public Subnet to Route Table: ", err)
	}

	fmt.Printf("[🐶] Successfully associated Public Subnet to Route Table\n")
	return ""
}
