package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a *CloudProvider) retrieveInternetGateway() string {
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

func (a *CloudProvider) createInternetGateway() string {
	svc := ec2.New(a.session)

	igw, err := svc.CreateInternetGateway(&ec2.CreateInternetGatewayInput{
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String("internet-gateway"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(a.awsConfig.VPC.InternetGateway.Name),
				},
			},
		}},
	})

	if err != nil {
		log.Fatal("[🐶] Error creating Internet Gateway: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", a.awsConfig.VPC.InternetGateway.Name, *igw.InternetGateway.InternetGatewayId)
	return *igw.InternetGateway.InternetGatewayId
}

func (a *CloudProvider) attachInternetGatewayToVpc(vpcId string, igwId string) string {
	svc := ec2.New(a.session)

	_, err := svc.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	})

	if err != nil {
		log.Fatal("[🐶] Error attaching Internet Gateway to VPC: ", err)
	}

	fmt.Printf("[🐶] Successfully attached %s to %s\n", a.awsConfig.VPC.InternetGateway.Name, a.awsConfig.VPC.Name)
	return ""
}

func (a *CloudProvider) addInternetGatewayToVpcPublicRouteTable(igwId string, publicRouteTableId string, publicSubnetId string) string {
	svc := ec2.New(a.session)

	_, err := svc.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(publicRouteTableId),
		GatewayId:            aws.String(igwId),
	})

	if err != nil {
		log.Fatal("[🐶] Error creating route: ", err)
	}

	fmt.Printf("[🐶] Successfully created Route for %s %s <-> %s association\n", a.awsConfig.VPC.Name, a.awsConfig.VPC.PublicRouteTable.Name, a.awsConfig.VPC.InternetGateway.Name)
	return ""
}

func (a *CloudProvider) associatePublicSubnetToPublicRouteTable(publicSubnetId string, publicRouteTableId string) string {
	svc := ec2.New(a.session)

	_, err := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		SubnetId:     aws.String(publicSubnetId),
		RouteTableId: aws.String(publicRouteTableId),
	})

	if err != nil {
		log.Fatal("[🐶] Error associating Public Subnet to Route Table: ", err)
	}

	fmt.Printf("[🐶] Successfully associated %s to %s\n", a.awsConfig.VPC.Subnets.PublicSubnetName, a.awsConfig.VPC.PublicRouteTable.Name)
	return ""
}
