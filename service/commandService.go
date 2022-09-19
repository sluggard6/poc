package service

import (
	"bytes"
	"os/exec"
)

type CommandService interface {
	Run(cmd string) (string, string, error)
}

var commandService = &commandImpl{"bash"}

func GetCommandService() CommandService {
	return commandService
}

type commandImpl struct {
	ShellToUse string
}

func (c *commandImpl) Run(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(c.ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
