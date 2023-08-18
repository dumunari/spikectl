package gcp

import (
	"context"
	"fmt"
	"log"

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

	_, err := g.client.Projects.Locations.Clusters.Create(createClusterRequest.Parent, createClusterRequest).Do()

	if err != nil {
		log.Fatal("[🐶] Error creating Cluster: ", err)
	}

	fmt.Printf("[🐶] Successfully created %s\n", g.gcpConfig.GKE.Name)
}

func (g *CloudProvider) retrieveKubeConfigInfo() config.KubeConfig {
	cluster := g.retrieveCluster()
	kubeConfig := config.KubeConfig{
		EndPoint: cluster.Endpoint,
		Token:    g.retrieveKubeToken(),
		CaFile:   g.retrieveCaFile(),
	}

	fmt.Println("[🐶] Kubeconfig successfully prepared")
	return kubeConfig
}

func (g *CloudProvider) retrieveKubeToken() string {
	kubeToken, err := g.credentials.TokenSource.Token()
	if err != nil {
		log.Fatal("[🐶] Error retrieving token: ", err)
	}

	return kubeToken.AccessToken
}

func (g *CloudProvider) retrieveCaFile() string {
	cluster := g.retrieveCluster()
	return utils.CreateTmpFile([]byte(cluster.MasterAuth.ClientCertificate))
}
