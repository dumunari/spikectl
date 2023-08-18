package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dumunari/spikectl/internal/cloud"
	"github.com/dumunari/spikectl/internal/config"
	"github.com/dumunari/spikectl/internal/core"
)

type Command interface {
	Execute() error
}

type InstallCommand struct {
	ConfigPath string
	Config     config.Spike
}

func (c *InstallCommand) Execute() error {
	spikeConfig, err := parseConfigFile(c.ConfigPath)
	if err != nil {
		return err
	}

	provider := cloud.NewCloudProvider(spikeConfig)

	kubeConfig := provider.InstantiateKubernetesCluster()

	core.InstallCoreComponents(spikeConfig, kubeConfig)

	return nil
}

func parseConfigFile(configPath string) (*config.Spike, error) {
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error while trying to open the file %s", configPath)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error while closing a file")
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error while trying to read the file %s", configPath)
		return nil, err
	}

	var cfg config.Spike
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("Error while trying to unmarshal the spike config json")
		return nil, err
	}

	return &cfg, nil
}
