package help

// Text returns the help text for the application
func Text() string {
	return `Creates directories and files in the current working directory, based on the structures defined in a scaff.json file (using the given variables).

SCAFF [commandname] [variablename]=[variablevalue]

SCAFF will work its way up the directory-tree, from the current working directory, searching for a scaff file that contains the requested command.
For full instructions on how to structure commands in a scaff.json file, visit https://github.com/M-Derbyshire/scaff

[commandname] - The name of the command (in a scaff.json file) that defines the files/directories to create.
[variablename]=[variablevalue] - Variables that are needed by the requested command can be defined in this format. See the below examples:

var1=myValue
var2="my longer value"

You can provide multiple variables in this way. If a variable is needed, but not provided, SCAFF will prompt you to provide it.

For full instructions on the use of SCAFF, visit https://github.com/M-Derbyshire/scaff`
}
