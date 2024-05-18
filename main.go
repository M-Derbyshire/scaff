package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/M-Derbyshire/scaff/commandfile"
	"github.com/M-Derbyshire/scaff/commandprocessing"
	"github.com/M-Derbyshire/scaff/helptext"
	"github.com/M-Derbyshire/scaff/uservariablemap"
)

func main() {
	configFileNameAndExt := "scaffconfig.json"
	args := os.Args[1:]
	workingDir, wdErr := os.Getwd()
	if wdErr != nil {
		panic(wdErr)
	}

	// Check a command name has been given (or a flag)
	if len(args) == 0 {
		fmt.Println("please provide the name of the command to process (or use '--help')")
		return
	}

	//Display help text
	if strings.EqualFold(args[0], "--help") {
		fmt.Println(helptext.GetHelpText())
		return
	}

	//Get the variables from the args
	var varMap map[string]string
	if len(args) > 1 { //first is the command name
		varMap = uservariablemap.GenerateVariableMap(args[1:])
	} else {
		varMap = make(map[string]string)
	}

	//Look for the command
	commandName := args[0]
	commandToProcess, commandTemplatePath, isFound := commandfile.FindCommand(commandName, configFileNameAndExt, workingDir)
	if !isFound {
		fmt.Println("unable to find the requested command ('" + commandName + "')")
		return
	}

	//Process command
	processingErr := commandprocessing.ProcessCommand(commandToProcess, workingDir, commandTemplatePath, varMap)
	if processingErr != nil {
		panic(processingErr)
	}
}
