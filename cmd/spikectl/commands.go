package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"spikectl/pkg/cloud"
	"spikectl/pkg/config"
)

type Command interface {
	Execute() error
}

type InstallCommand struct {
	ConfigPath string
	Config     config.SpikeConfig
}

func (c *InstallCommand) Execute() error {
	err := parseConfigFile(c.ConfigPath)
	if err != nil {
		return err
	}

	cloud.CreateCloudProvider(cfg)
}

func parseConfigFile(configPath string) (config.SpikeConfig, error) {
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
		fmt.Printf("Error while trying to read the file %s", c.ConfigPath)
		return err
	}

	var cfg config.SpikeConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("Error while trying to unmarshal the spike config json")
		return err
	}

	return nil
}
