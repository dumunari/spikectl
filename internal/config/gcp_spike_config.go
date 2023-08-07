package config

type GcpConfig struct {
	ProjectId string    `json:"project_id"`
	Zone      string    `json:"zone"`
	GKE       GKEConfig `json:"gke"`
	VPC       VPCConfig `json:"vpc"`
}

type GKEConfig struct {
	Name             string `json:"name"`
	Version          string `json:"version"`
	InitialNodeCount int64  `json:"initial_node_count"`
}
