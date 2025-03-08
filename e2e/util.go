package e2e_test

import (
	"bytes"
	"os/exec"
)

// This file contains util functions, used when constructing the E2E tests

// Runs the given shell command, from the given path, with the command arguments
func runShellCmd(startPath, command string, commandArgs ...string) (outStr string, errStr string, err error) {
	cmd := exec.Command(command, commandArgs...)
	cmd.Dir = startPath

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	cmdErr := cmd.Run()

	if cmdErr != nil {
		_, wasExitErr := cmdErr.(*exec.ExitError)

		if !wasExitErr {
			err = cmdErr
			return
		}
	}

	outStr = outb.String()
	errStr = errb.String()
	return
}
