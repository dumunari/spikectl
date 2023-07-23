package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

func (a CloudProvider) retrieveCluster() string {
	svc := eks.New(a.session)

	output, err := svc.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(a.awsConfig.EKS.Name),
	})

	if err != nil {
		//TODO: improve this error handling
		fmt.Printf("[🐶] Error describing Cluster: %s\n", err)
		return ""
	}

	fmt.Printf("[🐶] Found %s with Id: %s\n", a.awsConfig.EKS.Name, *output.Cluster.Id)
	return *output.Cluster.Id
}

func (a CloudProvider) createCluster(roleArn string, subnetIds ...string) {
	svc := eks.New(a.session)

	_, err := svc.CreateCluster(&eks.CreateClusterInput{
		Name: aws.String(a.awsConfig.EKS.Name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds: aws.StringSlice(subnetIds),
		},
		RoleArn: &roleArn,
	})

	if err != nil {
		log.Fatal("[🐶] Error creating Cluster: ", err)
	}

	//TODO: too messy
	fmt.Println("[🐶] Cluster creation requested, waiting for completion...")
	if err := svc.WaitUntilClusterActive(&eks.DescribeClusterInput{
		Name: aws.String(a.awsConfig.EKS.Name),
	}); err != nil {
		log.Fatal("[🐶] Error waiting for cluster creation: ", err)
	}
	fmt.Printf("[🐶] Successfully created %s\n", a.awsConfig.EKS.Name)
}
