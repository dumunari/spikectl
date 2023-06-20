package main

type Command interface {
	Execute()
}

type InstallCommand struct {
	ConfigPath string
}

func (c *InstallCommand) Execute() {
	panic("implement me")
}
