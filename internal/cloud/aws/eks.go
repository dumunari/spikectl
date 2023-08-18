package aws

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/dumunari/spikectl/internal/config"
	"github.com/dumunari/spikectl/internal/utils"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

func (a *CloudProvider) retrieveCluster() eks.Cluster {
	svc := eks.New(a.session)

	output, err := svc.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(a.awsConfig.EKS.Name),
	})

	if err != nil {
		//TODO: improve this error handling
		fmt.Printf("[🐶] Error describing Cluster: %s\n", err)
		return eks.Cluster{}
	}

	fmt.Printf("[🐶] Found %s with Arn: %s\n", a.awsConfig.EKS.Name, *output.Cluster.Arn)
	return *output.Cluster
}

func (a *CloudProvider) createCluster(roleArn string, subnetIds ...string) eks.Cluster {
	svc := eks.New(a.session)

	output, err := svc.CreateCluster(&eks.CreateClusterInput{
		Name: aws.String(a.awsConfig.EKS.Name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds:             aws.StringSlice(subnetIds),
			EndpointPrivateAccess: aws.Bool(true),
			EndpointPublicAccess:  aws.Bool(true),
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
	return *output.Cluster
}

func (a *CloudProvider) retrieveKubeConfigInfo(cluster eks.Cluster) config.KubeConfig {
	kubeConfig := config.KubeConfig{
		EndPoint: *cluster.Endpoint,
		Token:    a.retrieveKubeToken(cluster),
		CaFile:   a.retrieveCaFile(cluster),
	}

	fmt.Println("[🐶] Kubeconfig successfully prepared")
	return kubeConfig
}

func (a *CloudProvider) retrieveKubeToken(cluster eks.Cluster) string {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		log.Fatal("[🐶] Error on token generator: ", err)
	}

	tok, err := gen.GetWithOptions(&token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
		Session:   a.session,
	})
	if err != nil {
		log.Fatal("[🐶] Error generating token: ", err)
	}

	return tok.Token
}

func (a *CloudProvider) retrieveCaFile(cluster eks.Cluster) string {
	caData, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		log.Fatal("[🐶] Error retrieving CA Data: ", err)
	}

	return utils.CreateTmpFile(caData)
}
