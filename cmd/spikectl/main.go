package main

import (
	"flag"
	"log"
	"os"
	"spikectl/internal/commands"
)

func main() {
	var command commands.Command
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			command = parseInstallCommand()
		}

		err := command.Execute()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No command provided")
	}
}

func parseInstallCommand() commands.Command {
	var command commands.InstallCommand
	installFlagSet := flag.NewFlagSet("install", flag.ExitOnError)
	installFlagSet.StringVar(&command.ConfigPath, "configPath", "./.spikecfg", "The path to the spike config file. If none is provided, the default path will be used (./.spikecfg)")

	err := installFlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Fatal("Error while parsing command arguments")
	}

	return &command
}
