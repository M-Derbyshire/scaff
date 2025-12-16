package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/M-Derbyshire/scaff/command"
	"github.com/M-Derbyshire/scaff/help"
	"github.com/M-Derbyshire/scaff/variable"
)

func main() {
	scaffFileNameAndExt := "scaff.json"
	args := os.Args[1:]
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Check a command name has been given (or a flag)
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "please provide the name of the command to process (or use '--help')")
		os.Exit(1)
	}

	//Display help text
	if strings.EqualFold(args[0], "--help") {
		fmt.Println(help.Text())
		return
	}

	//Get the variables from the args
	var varMap map[string]string
	if len(args) > 1 { //first is the command name
		varMap = variable.Map(args[1:])
	} else {
		varMap = make(map[string]string)
	}

	//Look for the command
	commandName := args[0]
	commandToProcess, fullTemplatePath, isFound, err := command.Find(commandName, scaffFileNameAndExt, workingDir)
	if err != nil {
		panic(err)
	}
	if !isFound {
		fmt.Fprintln(os.Stderr, "unable to find the requested command ('"+commandName+"')")
		os.Exit(3)
	}

	// Confirm that no files/directories in the command already exist
	existingPaths, err := command.IdentifyExistingPaths(commandToProcess, workingDir, varMap)
	if err != nil {
		panic(err)
	}
	if len(existingPaths) > 0 {
		for _, path := range existingPaths {
			fmt.Fprintln(os.Stderr, "path already exists:", path)
		}

		os.Exit(4)
	}

	//Process command
	err = command.Process(commandToProcess, workingDir, fullTemplatePath, varMap)
	if err != nil {
		panic(err)
	}
}
