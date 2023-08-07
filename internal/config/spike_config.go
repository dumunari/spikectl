package config

type SpikeConfig struct {
	IDP IDPConfig `json:"idp"`
}

type IDPConfig struct {
	CloudProvider CloudProvider `json:"cloud_provider"`
	AzureConfig   AzureConfig   `json:"azure"`
	AwsConfig     AwsConfig     `json:"aws"`
	GcpConfig     GcpConfig     `json:"gcp"`
}
