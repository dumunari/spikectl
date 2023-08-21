package config

type Spike struct {
	Spike SpikeConfig `json:"spike"`
}

type SpikeConfig struct {
	CloudProvider  CloudProvider   `json:"cloud_provider"`
	AzureConfig    AzureConfig     `json:"azure"`
	AwsConfig      AwsConfig       `json:"aws"`
	GcpConfig      GcpConfig       `json:"gcp"`
	CoreComponents []CoreComponent `json:"core"`
}

type CoreComponent struct {
	ReleaseName  string `json:"release_name"`
	Chart        string `json:"chart"`
	Namespace    string `json:"namespace"`
	Repository   string `json:"repository"`
	ChartVersion string `json:"chart_version"`
	ImageName    string `json:"image_name,omitempty"`
	ImageVersion string `json:"image_version,omitempty"`
}
