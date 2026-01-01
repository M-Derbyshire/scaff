package e2e

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/creack/pty"
)

// This file contains util functions, used when constructing the E2E tests

// Runs the given shell command, from the given path, with the command arguments.
// The inputs slice is used to respond when the app prompts the user for variable values.
func runShellCmd(startPath, command string, inputs []string, commandArgs ...string) (outStr string, errStr string, err error) {
	operatingSystem := runtime.GOOS

	cmd := exec.Command(command, commandArgs...)
	cmd.Dir = startPath

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if operatingSystem == "windows" {
		stdinReader, stdinWriter := io.Pipe()
		cmd.Stdin = stdinReader

		err = cmd.Start()
		if err != nil {
			return
		}

		// Write inputs sequentially, in case the app prompts for input
		go func() {
			for _, input := range inputs {
				_, writeErr := fmt.Fprintln(stdinWriter, input)
				if writeErr != nil {
					return
				}
				time.Sleep(100 * time.Millisecond) // Gives the app some time to prompt for the next input
			}

			stdinWriter.Close()
		}()
	} else {
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return "", "", err
		}
		defer ptmx.Close()

		var output bytes.Buffer

		// Capture output
		go func() {
			io.Copy(&output, ptmx)
		}()

		// Feed inputs
		go func() {
			for _, input := range inputs {
				fmt.Fprintln(ptmx, input)
			}
		}()
	}

	cmdErr := cmd.Wait()

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
// The inputs slice is used to respond when the app prompts the user for variable values.
// The args are any other arguments to pass to the app.
func runScaffoldCommand(commandName string, inputs []string, args ...string) error {
	allCmdArgs := append([]string{commandName}, args...)

	_, errStr, err := runShellCmd(scaffoldRunPath, "./scaff", inputs, allCmdArgs...)
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
	outStr, _, err := runShellCmd(".", "diff", []string{}, "-r", scaffoldRunPath, fmt.Sprintf("./expected/%v", commandName))
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
