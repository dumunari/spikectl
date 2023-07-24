package config

import (
	"testing"
)

func TestCloudProvider_MarshalJSON(t *testing.T) {
	testParameter := map[CloudProvider]string{
		AWS:   "\"AWS\"",
		AZURE: "\"Azure\"",
		GCP: "\"GCP\"",
	}

	for provider, expectedResult := range testParameter {
		got, err := provider.MarshalJSON()

		if err != nil {
			t.Errorf("Error while trying to marshal cloud provider to json")
		}

		if string(got) != expectedResult {
			t.Errorf("got %s, wanted %s", string(got), expectedResult)
		}
	}
}

func TestCloudProvider_UnmarshalJSON(t *testing.T) {
	testParameters := map[string]CloudProvider{
		"\"AWS\"":   AWS,
		"\"Azure\"": AZURE,
		"\"GCP\"":   GCP,
	}

	for value, expectedProvider := range testParameters {
		var cloudProvider CloudProvider
		err := cloudProvider.UnmarshalJSON([]byte(value))
		if err != nil {
			t.Errorf("Error while trying to unmarshal json")
		}

		if cloudProvider != expectedProvider {
			t.Errorf("Expected cloudprovider %q, got %q", expectedProvider, cloudProvider)
		}
	}
}
