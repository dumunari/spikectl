package config

type AwsConfig struct {
	Profile string    `json:"profile"`
	Region  string    `json:"region"`
	VPC     VPCConfig `json:"vpc"`
	EKS     EKSConfig `json:"eks"`
}

type VPCConfig struct {
	Name    string       `json:"name"`
	CIDR    string       `json:"cidr"`
	Subnets SubnetConfig `json:"subnets"`
}

type SubnetConfig struct {
	PublicSubnetName  string `json:"public_subnet_name"`
	PublicSubnetCidr  string `json:"public_subnet_cidr"`
	PublicSubnetAz    string `json:"public_subnet_az"`
	PrivateSubnetName string `json:"private_subnet_name"`
	PrivateSubnetCidr string `json:"private_subnet_cidr"`
	PrivateSubnetAz   string `json:"private_subnet_az"`
}

type EKSConfig struct {
	Name string `json:"name"`
}
