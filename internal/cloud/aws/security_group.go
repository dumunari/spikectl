package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (a CloudProvider) retrieveSecurityGroup(vpcId string) string {
	svc := ec2.New(a.session)

	output, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(a.awsConfig.EKS.SecurityGroup.Name)},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	})

	if err != nil {
		fmt.Printf("[🐶] Error describing Security Groups: %s\n", err)
		return ""
	}

	if len(output.SecurityGroups) == 0 {
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", *output.SecurityGroups[0].GroupName, *output.SecurityGroups[0].GroupId)

	return *output.SecurityGroups[0].GroupId
}

func (a CloudProvider) createSecurityGroup(vpcId string) string {
	svc := ec2.New(a.session)

	sg, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(a.awsConfig.EKS.SecurityGroup.Name),
		VpcId:       aws.String(vpcId),
		Description: aws.String("Security Group For IDP"),
	})

	if err != nil {
		log.Fatal("[🐶] Error creating Security Group: ", err)
	}

	fmt.Printf("[🐶] %s Successfully created: %s\n", a.awsConfig.EKS.SecurityGroup.Name, *sg.GroupId)

	return *sg.GroupId
}
