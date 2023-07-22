package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/dumunari/spikectl/internal/config"
)

type CloudProvider struct {
	awsConfig config.AwsConfig
	session   *session.Session
}

func (a *CloudProvider) InstantiateKubernetesCluster() error {
	vpcId := a.retrieveVpc()
	if vpcId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Name)
		vpcId = a.createVpc()
	}

	publicSubnetId := a.retrieveSubnet(a.awsConfig.VPC.Subnets.PublicSubnetName)
	if publicSubnetId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Subnets.PublicSubnetName)
		a.createSubnet(&vpcId, a.awsConfig.VPC.Subnets.PublicSubnetName, a.awsConfig.VPC.Subnets.PublicSubnetCidr)
	}

	privateSubnetId := a.retrieveSubnet(a.awsConfig.VPC.Subnets.PrivateSubnetName)
	if privateSubnetId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Subnets.PrivateSubnetName)
		a.createSubnet(&vpcId, a.awsConfig.VPC.Subnets.PrivateSubnetName, a.awsConfig.VPC.Subnets.PrivateSubnetCidr)
	}

	return nil
}

func NewAwsCloudProvider(config *config.SpikeConfig) *CloudProvider {
	awsProvider := CloudProvider{}
	awsProvider.awsConfig = config.IDP.AwsConfig

	if len(awsProvider.awsConfig.Region) == 0 {
		log.Fatal("[🐶] No AWS Region provided.")
	}

	if len(awsProvider.awsConfig.Profile) == 0 {
		log.Fatal("[🐶] No AWS Profile provided.")
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(awsProvider.awsConfig.Region),
		},
		Profile: awsProvider.awsConfig.Profile,
	})

	if err != nil {
		log.Fatal("[🐶] Error creating AWS Config")
	}

	awsProvider.session = sess

	return &awsProvider
}
