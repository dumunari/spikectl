package gcp

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/container/v1"
)

func (a *CloudProvider) retrieveCluster() string {
	ctx := context.Background()

	parent := fmt.Sprintf("projects/%s/locations/%s", a.gcpConfig.ProjectId, a.gcpConfig.Zone)

	resp, err := a.client.Projects.Locations.Clusters.List(parent).Context(ctx).Do()

	if err != nil {
		log.Fatal("[🐶] Error listing GKEs: ", err)
	}

	if len(resp.Clusters) == 0 {
		return ""
	}

	for _, cluster := range resp.Clusters {
		if cluster.Name == a.gcpConfig.GKE.Name {
			fmt.Printf("[🐶] Found %s with Id: %s\n", a.gcpConfig.GKE.Name, cluster.Id)
			return cluster.SelfLink
		}
	}

	return ""
}

func (a *CloudProvider) createCluster(vpcLink, publicSubnetLink string) {
	cluster := &container.Cluster{
		Name:                  a.gcpConfig.GKE.Name,
		Network:               vpcLink,
		Subnetwork:            publicSubnetLink,
		Zone:                  a.gcpConfig.Zone,
		InitialNodeCount:      a.gcpConfig.GKE.InitialNodeCount,
		InitialClusterVersion: a.gcpConfig.GKE.Version,
	}

	createClusterRequest := &container.CreateClusterRequest{
		Cluster:   cluster,
		Parent:    fmt.Sprintf("projects/%s/locations/%s", a.gcpConfig.ProjectId, a.gcpConfig.Zone),
		ProjectId: a.gcpConfig.ProjectId,
		Zone:      a.gcpConfig.Zone,
	}

	_, err := a.client.Projects.Locations.Clusters.Create(createClusterRequest.Parent, createClusterRequest).Do()

	if err != nil {
		log.Fatal("[🐶] Error creating Cluster: ", err)
	}

	fmt.Printf("[🐶] Successfully created %s\n", a.gcpConfig.GKE.Name)
}
