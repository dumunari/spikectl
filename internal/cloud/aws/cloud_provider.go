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

// TODO: improve this method
func (a *CloudProvider) InstantiateKubernetesCluster() error {
	fmt.Println("[🐶] Checking network resources...")

	vpcId := a.retrieveVpc()
	if vpcId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Name)
		vpcId = a.createVpc()
	}
	mainRouteTableId := a.retrieveMainRouteTable(vpcId)

	publicRouteTableId := a.retrievePublicRouteTable(vpcId)
	if publicRouteTableId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Name)
		publicRouteTableId = a.createPublicRouteTable(vpcId)
	}

	igwId := a.retrieveInternetGateway()
	if igwId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.InternetGateway.Name)
		igwId = a.createInternetGateway()
		a.attachInternetGateway(vpcId, igwId)
	}

	publicSubnetId := a.retrieveSubnet(a.awsConfig.VPC.Subnets.PublicSubnetName)
	if publicSubnetId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Subnets.PublicSubnetName)
		publicSubnetId = a.createSubnet(&vpcId, a.awsConfig.VPC.Subnets.PublicSubnetName, a.awsConfig.VPC.Subnets.PublicSubnetCidr, a.awsConfig.VPC.Subnets.PublicSubnetAz)
		a.addPublicIpAutoAssignToSubnet(publicSubnetId)
		a.addInternetGatewayToVpcPublicRouteTable(igwId, publicRouteTableId, publicSubnetId)
		a.associatePublicSubnetToPublicRouteTable(publicSubnetId, publicRouteTableId)
	}

	ngwId := a.retrieveNatGateway()
	if ngwId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.NatGateway.Name)
		ngwId := a.createNatGateway(publicSubnetId)
		a.addNatGatewayToVpcMainRouteTable(mainRouteTableId, ngwId)
	}

	primaryPrivateSubnetId := a.retrieveSubnet(a.awsConfig.VPC.Subnets.PrimaryPrivateSubnetName)
	if primaryPrivateSubnetId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Subnets.PrimaryPrivateSubnetName)
		a.createSubnet(&vpcId, a.awsConfig.VPC.Subnets.PrimaryPrivateSubnetName, a.awsConfig.VPC.Subnets.PrimaryPrivateSubnetCidr, a.awsConfig.VPC.Subnets.PrimaryPrivateSubnetAz)
	}

	secondaryPrivateSubnetId := a.retrieveSubnet(a.awsConfig.VPC.Subnets.SecondaryPrivateSubnetName)
	if secondaryPrivateSubnetId == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.VPC.Subnets.SecondaryPrivateSubnetName)
		a.createSubnet(&vpcId, a.awsConfig.VPC.Subnets.SecondaryPrivateSubnetName, a.awsConfig.VPC.Subnets.SecondaryPrivateSubnetCidr, a.awsConfig.VPC.Subnets.SecondaryPrivateSubnetAz)
	}

	fmt.Println("[🐶] Checking Kubernetes Cluster...")

	eksCluster := a.retrieveCluster()
	if len(eksCluster) == 0 {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.EKS.Name)

		//TODO: improve role and attachRole actions
		iamClusterRole := a.retrieveClusterRole()
		if iamClusterRole == "" {
			fmt.Println("[🐶] No eksClusterRole found, creating one...")
			iamClusterRole = a.createClusterRole()
		}
		a.attachClusterRolePolicy()

		a.createCluster(iamClusterRole, publicSubnetId, primaryPrivateSubnetId, secondaryPrivateSubnetId)
	}

	nodeGroup := a.retrieveNodeGroup()
	if nodeGroup == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.awsConfig.EKS.NodeGroup.Name)
		iamNodeRole := a.retrieveNodeRole()
		if iamNodeRole == "" {
			fmt.Println("[🐶] No eksNodeRole found, creating one...")
			iamNodeRole = a.createNodeRole()
		}
		a.attachNodeRolePolicy()

		a.createNodeGroup(iamNodeRole, publicSubnetId, primaryPrivateSubnetId, secondaryPrivateSubnetId)
	}

	return nil
}

func NewAwsCloudProvider(config *config.Spike) *CloudProvider {
	awsProvider := CloudProvider{}
	awsProvider.awsConfig = config.Spike.AwsConfig

	fmt.Println(config.Spike.AwsConfig)

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
