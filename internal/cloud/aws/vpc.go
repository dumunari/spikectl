package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a CloudProvider) retrieveVpc() string {
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

func (a CloudProvider) createVpc() string {
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

	fmt.Printf("[🐶] %s Succesfully created: %s\n", a.awsConfig.VPC.Name, *vpc.Vpc.VpcId)

	return *vpc.Vpc.VpcId
}
