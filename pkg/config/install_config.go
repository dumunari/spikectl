package config

import (
	"encoding/json"
	"fmt"
)

type CloudProvider uint8

const (
	AWS CloudProvider = iota + 1
	AZURE
)

func (c *CloudProvider) String() string {
	switch *c {
	case AWS:
		return "AWS"
	case AZURE:
		return "Azure"
	}
	return "unknown"
}
func (c *CloudProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
func (c *CloudProvider) UnmarshalJSON(data []byte) (err error) {
	var cloudprovider string
	if err := json.Unmarshal(data, &cloudprovider); err != nil {
		return err
	}
	if *c, err = ParseCloudProvider(cloudprovider); err != nil {
		return err
	}
	return nil
}

func ParseCloudProvider(s string) (CloudProvider, error) {
	switch s {
	case "AWS":
		return AWS, nil
	case "Azure":
		return AZURE, nil
	}
	return CloudProvider(0), fmt.Errorf("%q is not a supported cloud platform", s)
}

type IDPConfig struct {
	CloudProvider string `json:"cloud_provider"`
}

type SpikeConfig struct {
	IDP IDPConfig `json:"idp"`
}
