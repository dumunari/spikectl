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
	}

	for value, expectedProvider := range testParameters {
		js := fmt.Sprintf("{\"idp\": {\"cloud_provider\": \"%s\"}}", value)
		var spikeConfig SpikeConfig
		err := json.Unmarshal([]byte(js), &spikeConfig)

		if err != nil {
			t.Errorf("Error while trying to unmarshal an spikeconfig json")
		}

		if expectedProvider != spikeConfig.IDP.CloudProvider {
			t.Errorf("Expected clooud_provider %q, got %q", expectedProvider, spikeConfig.IDP.CloudProvider)
		}
	}
}
