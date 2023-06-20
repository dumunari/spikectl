package main

import (
	"flag"
	"os"
	"spikectl/cmd/spikectl/commands"
)

func main() {
	var command commands.Command
	if os.Args[1] == "install" {
		command = parseInstallCommand()
	}

	command.Execute()
}

func parseInstallCommand() commands.Command {
	var command commands.InstallCommand
	installFlagSet := flag.NewFlagSet("install", flag.ExitOnError)
	installFlagSet.StringVar(&command.ConfigPath, "configPath", "./.spikecfg", "The path to the spike config file. If none is provided, the default path will be used (./spikecfg)")

	return &command
}
