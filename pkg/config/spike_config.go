package config

type IDPConfig struct {
	CloudProvider CloudProvider `json:"cloud_provider"`
}

type SpikeConfig struct {
	IDP IDPConfig `json:"idp"`
}
