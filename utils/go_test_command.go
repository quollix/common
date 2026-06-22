package utils

import "github.com/quollix/taskrunner"

type GoTestCommand struct {
	directory string
	buildTag  string
	runFilter string
	env       map[string]string
}

func GoTest(directory string) *GoTestCommand {
	return &GoTestCommand{
		directory: directory,
		env:       make(map[string]string),
	}
}

func (c *GoTestCommand) Tag(buildTag string) *GoTestCommand {
	c.buildTag = buildTag
	return c
}

func (c *GoTestCommand) Filter(runFilter string) *GoTestCommand {
	c.runFilter = runFilter
	return c
}

func (c *GoTestCommand) Env(key, value string) *GoTestCommand {
	c.env[key] = value
	return c
}

func (c *GoTestCommand) Run(tr *taskrunner.TaskRunner) {
	command := tr.Cmd().Dir(c.directory)
	for key, value := range c.env {
		command.Env(key, value)
	}
	command.Run("%s", c.buildCommand())
}

func (c *GoTestCommand) buildCommand() string {
	command := "go test -p 1 -v -count=1 -failfast"
	if c.buildTag != "" {
		command += " -tags=" + c.buildTag
	}
	if c.runFilter != "" {
		command += " -run '" + c.runFilter + "'"
	}
	command += " ./..."
	return command
}
