package config

type AwsConfig struct {
	Profile string    `json:"profile"`
	Region  string    `json:"region"`
	VPC     VPCConfig `json:"vpc"`
	EKS     EKSConfig `json:"eks"`
}

type VPCConfig struct {
	Name             string                `json:"name"`
	CIDR             string                `json:"cidr"`
	Subnets          SubnetConfig          `json:"subnets"`
	InternetGateway  InternetGatewayConfig `json:"internet_gateway"`
	NatGateway       NatGatewayConfig      `json:"nat_gateway"`
	PublicRouteTable RouteTableConfig      `json:"public_route_table"`
}

type SubnetConfig struct {
	PublicSubnetName           string `json:"public_subnet_name"`
	PublicSubnetCidr           string `json:"public_subnet_cidr"`
	PublicSubnetAz             string `json:"public_subnet_az"`
	PrimaryPrivateSubnetName   string `json:"primary_private_subnet_name"`
	PrimaryPrivateSubnetCidr   string `json:"primary_private_subnet_cidr"`
	PrimaryPrivateSubnetAz     string `json:"primary_private_subnet_az"`
	SecondaryPrivateSubnetName string `json:"secondary_private_subnet_name"`
	SecondaryPrivateSubnetCidr string `json:"secondary_private_subnet_cidr"`
	SecondaryPrivateSubnetAz   string `json:"secondary_private_subnet_az"`
}

type InternetGatewayConfig struct {
	Name string `json:"name"`
}

type NatGatewayConfig struct {
	Name string `json:"name"`
}

type RouteTableConfig struct {
	Name string `json:"name"`
}

type EKSConfig struct {
	Name      string          `json:"name"`
	NodeGroup NodeGroupConfig `json:"node_group"`
}

type NodeGroupConfig struct {
	Name         string        `json:"name"`
	InstanceType []string      `json:"instance_type"`
	Scaling      ScalingConfig `json:"scaling"`
}

type ScalingConfig struct {
	MaxSize     int64 `json:"max_size"`
	MinSize     int64 `json:"min_size"`
	DesiredSize int64 `json:"desired_size"`
}
