package config

type Spike struct {
	Spike SpikeConfig `json:"spike"`
}

type SpikeConfig struct {
	CloudProvider CloudProvider   `json:"cloud_provider"`
	AzureConfig   AzureConfig     `json:"azure"`
	AwsConfig     AwsConfig       `json:"aws"`
	GcpConfig     GcpConfig       `json:"gcp"`
	CoreConfig    []CoreComponent `json:"core"`
}

type CoreComponent struct {
	ReleaseName string `json:"release_name"`
	Chart       string `json:"chart"`
	Namespace   string `json:"namespace"`
	Repository  string `json:"repository"`
	Version     string `json:"version"`
}
