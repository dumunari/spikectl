package gcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dumunari/spikectl/internal/config"
	"github.com/dumunari/spikectl/internal/utils"
	"google.golang.org/api/container/v1"
)

func (g *CloudProvider) retrieveCluster() container.Cluster {
	ctx := context.Background()

	parent := fmt.Sprintf("projects/%s/locations/%s", g.gcpConfig.ProjectId, g.gcpConfig.Zone)

	resp, err := g.client.Projects.Locations.Clusters.List(parent).Context(ctx).Do()

	if err != nil {
		log.Fatal("[🐶] Error listing GKEs: ", err)
	}

	if len(resp.Clusters) == 0 {
		return container.Cluster{}
	}

	for _, cluster := range resp.Clusters {
		if cluster.Name == g.gcpConfig.GKE.Name {
			fmt.Printf("[🐶] Found %s with Id: %s\n", g.gcpConfig.GKE.Name, cluster.Id)
			return *cluster
		}
	}

	return container.Cluster{}
}

func (g *CloudProvider) createCluster(vpcLink, publicSubnetLink string) {
	ctx := context.Background()

	cluster := &container.Cluster{
		Name:                  g.gcpConfig.GKE.Name,
		Network:               vpcLink,
		Subnetwork:            publicSubnetLink,
		Zone:                  g.gcpConfig.Zone,
		InitialNodeCount:      g.gcpConfig.GKE.InitialNodeCount,
		InitialClusterVersion: g.gcpConfig.GKE.Version,
	}

	createClusterRequest := &container.CreateClusterRequest{
		Cluster:   cluster,
		Parent:    fmt.Sprintf("projects/%s/locations/%s", g.gcpConfig.ProjectId, g.gcpConfig.Zone),
		ProjectId: g.gcpConfig.ProjectId,
		Zone:      g.gcpConfig.Zone,
	}

	op, err := g.client.Projects.Locations.Clusters.Create(createClusterRequest.Parent, createClusterRequest).Do()

	if err != nil {
		log.Fatal("[🐶] Error creating Cluster: ", err)
	}

	err = WaitForOperation(ctx, op.Name, g.client, g.gcpConfig.ProjectId, g.gcpConfig.Zone)

	if err != nil {
		log.Fatal("[🐶] Error waiting for operation: ", err)
	}

	fmt.Printf("[🐶] Successfully created %s\n", g.gcpConfig.GKE.Name)
}

func WaitForOperation(ctx context.Context, opName string, service *container.Service, projectID, zone string) error {
	for {

		operation, err := service.Projects.Zones.Operations.Get(projectID, zone, opName).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("[🐶] Error waiting for operation: %v", err)
		}

		if operation.Status == "DONE" {
			if operation.StatusMessage != "" {
				return fmt.Errorf(operation.StatusMessage)
			}
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (g *CloudProvider) retrieveKubeConfigInfo() config.KubeConfig {
	cluster := g.retrieveCluster()
	kubeConfig := config.KubeConfig{
		EndPoint: cluster.Endpoint,
		Token:    g.retrieveKubeToken(cluster),
		CaFile:   g.retrieveCaFile(cluster),
	}

	fmt.Println("[🐶] Kubeconfig successfully prepared")
	return kubeConfig
}

func (g *CloudProvider) retrieveKubeToken(cluster container.Cluster) string {
	kubeToken, err := g.credentials.TokenSource.Token()
	if err != nil {
		log.Fatal("[🐶] Error retrieving token: ", err)
	}

	return kubeToken.AccessToken
}

func (g *CloudProvider) retrieveCaFile(cluster container.Cluster) string {
	return utils.CreateTmpFile([]byte(cluster.MasterAuth.ClientCertificate))
}
