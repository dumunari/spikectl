package main

import (
	"flag"
	"os"
)

func main() {
	var command Command
	if os.Args[1] == "install" {
		command = parseInstallCommand()
	}

	command.Execute()
}

func parseInstallCommand() Command {
	var command InstallCommand
	installFlagSet := flag.NewFlagSet("install", flag.ExitOnError)
	installFlagSet.StringVar(&command.ConfigPath, "configPath", "./.spikecfg", "The path to the spike config file. If none is provided, the default path will be used (./spikecfg)")

	return &command
}
