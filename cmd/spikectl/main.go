package main

import (
	"flag"
	"log"
	"os"
	"spikectl/internal/commands"
)

func main() {
	var command commands.Command
	if os.Args[1] == "install" {
		command = parseInstallCommand()
	}

	err := command.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func parseInstallCommand() commands.Command {
	var command commands.InstallCommand
	installFlagSet := flag.NewFlagSet("install", flag.ExitOnError)
	installFlagSet.StringVar(&command.ConfigPath, "configPath", "./.spikecfg", "The path to the spike config file. If none is provided, the default path will be used (./spikecfg)")

	return &command
}
