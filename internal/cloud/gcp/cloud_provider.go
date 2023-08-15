package gcp

import (
	"context"
	"fmt"
	"log"

	"github.com/dumunari/spikectl/internal/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
)

type CloudProvider struct {
	gcpConfig   config.GcpConfig
	client      *container.Service
	credentials *google.Credentials
}

func NewGcpCloudProvider(config *config.Spike) *CloudProvider {
	gcpProvider := CloudProvider{}
	gcpProvider.gcpConfig = config.Spike.GcpConfig

	if len(gcpProvider.gcpConfig.ProjectId) == 0 {
		log.Fatal("[🐶] No GCP Project ID provided.")
	}

	if len(gcpProvider.gcpConfig.Zone) == 0 {
		log.Fatal("[🐶] No GCP Zone provided.")
	}

	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, container.CloudPlatformScope)
	if err != nil {
		log.Fatal(fmt.Sprintf("[🐶] Error getting credentials: %v", err))
	}

	client, err := container.NewService(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("[🐶] Error creating GCP client: %v", err))
	}

	gcpProvider.client = client
	gcpProvider.credentials = creds

	return &gcpProvider
}

func (a *CloudProvider) InstantiateKubernetesCluster() error {
	vpcLink := a.getOrCreateVpc(a.gcpConfig.VPC.Name)

	publicSubnetId := a.getOrCreateSubnet(vpcLink, a.gcpConfig.VPC.Subnets.PublicSubnetName, a.gcpConfig.VPC.Subnets.PublicSubnetAz, a.gcpConfig.VPC.Subnets.PublicSubnetCidr)

	cluster := a.retrieveCluster()
	if cluster == "" {
		fmt.Printf("[🐶] No %s found, creating one...\n", a.gcpConfig.GKE.Name)
		a.createCluster(vpcLink, publicSubnetId)
	}

	return nil
}
