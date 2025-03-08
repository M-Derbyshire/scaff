package e2e_test

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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

var scaffoldRunPath = "./environment/child_dir/grandchild_dir/" // This is the location to run any scaff commands from

// Runs a scaffold command in the environment directory
// The commandName is the name of the command in the scaff.json file.
// The args are any other arguments to pass to the app.
func runScaffoldCommand(commandName string, args ...string) error {
	allCmdArgs := append([]string{commandName}, args...)

	_, errStr, err := runShellCmd(scaffoldRunPath, "./scaff", allCmdArgs...)
	if err != nil {
		return err
	}

	if len(errStr) > 0 {
		return errors.New(errStr)
	}

	return nil
}

// Uses the DIFF command to compare the results of a scaffold action with the expected output
func diffScaffoldCommand(commandName string) ([]string, error) {
	outStr, _, err := runShellCmd(".", "diff", "-r", scaffoldRunPath, fmt.Sprintf("./expected/%v", commandName))
	if err != nil {
		return []string{}, err
	}

	allDiffs := strings.Split(outStr, "\n")

	return getRelevantDiffResults(allDiffs), nil
}

// Takes a slice of diff results, and filters to just the relevant ones
func getRelevantDiffResults(allDiffs []string) []string {
	var relevantDiffs []string

	for _, diff := range allDiffs {
		if !diffCanBeIgnored(diff) {
			relevantDiffs = append(relevantDiffs, diff)
		}
	}

	return relevantDiffs
}

// Determines if a diff result isn't relevant
func diffCanBeIgnored(diff string) bool {
	// These are the file/dir names we can ignore
	namesToIgnore := []string{
		"my_templates",
		"scaff_files",
		"scaff",
		"scaff.exe",
		"scaff.json",
	}

	if len(diff) == 0 {
		return true
	}

	for _, ignoredName := range namesToIgnore {
		if strings.HasSuffix(diff, ignoredName) {
			return true
		}
	}

	return false
}
