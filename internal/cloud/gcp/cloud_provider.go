package azure

import (
	"log"
	"context"
	"fmt"

	"google.golang.org/api/container/v1"
	"golang.org/x/oauth2/google"
	"github.com/dumunari/spikectl/internal/config"
)

type CloudProvider struct {
	gcpConfig 	config.GcpConfig
	client    	*container.Service
	credentials *google.Credentials
}

func NewGcpCloudProvider(config *config.SpikeConfig) *CloudProvider {
	gcpProvider := CloudProvider{}
	gcpProvider.gcpConfig = config.IDP.GcpConfig

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
	cluster := &container.Cluster{
		Name: a.gcpConfig.GKE.Name,
		Zone: a.gcpConfig.Zone,
		InitialNodeCount:  a.gcpConfig.GKE.InitialNodeCount,
		InitialClusterVersion: a.gcpConfig.GKE.Version,
	}

	createClusterRequest := &container.CreateClusterRequest{
		Cluster: cluster,
		Parent: fmt.Sprintf("projects/%s/locations/%s", a.gcpConfig.ProjectId, a.gcpConfig.Zone),
		ProjectId: a.gcpConfig.ProjectId,
		Zone: a.gcpConfig.Zone,
	}


	_, err := a.client.Projects.Locations.Clusters.Create(createClusterRequest.Parent, createClusterRequest).Do()

	if err != nil {
		log.Fatal(fmt.Sprintf("[🐶] Error creating GKE cluster: %v", err))
	}

	return nil
}
