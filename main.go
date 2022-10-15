package main

import (
	"fmt"
	"os"

	"github.com/M-Derbyshire/scaff/uservariablemap"
)

func main() {
	userVariables := uservariablemap.GenerateVariableMap(os.Args[2:]) //Get the arguments after the app name, and the user's called command
	fmt.Println(userVariables)
}
