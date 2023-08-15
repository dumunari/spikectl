package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseCloudProvider(t *testing.T) {
	testParameters := map[string]CloudProvider{
		"Azure": AZURE,
		"AWS":   AWS,
		"GCP":   GCP,
	}

	for value, expectedProvider := range testParameters {
		js := fmt.Sprintf("{\"spike\": {\"cloud_provider\": \"%s\"}}", value)
		var spikeConfig Spike
		err := json.Unmarshal([]byte(js), &spikeConfig)

		if err != nil {
			t.Errorf("Error while trying to unmarshal an spikeconfig json")
		}

		if expectedProvider != spikeConfig.Spike.CloudProvider {
			t.Errorf("Expected cloud_provider %q, got %q", expectedProvider, spikeConfig.Spike.CloudProvider)
		}
	}
}
