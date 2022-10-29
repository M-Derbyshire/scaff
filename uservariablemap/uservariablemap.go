package uservariablemap

import "strings"

// GenerateVariableMap when given a slice of string arguments, this will convert any argument
// that contains an "=" symbol into a key/value pair in the returned map
func GenerateVariableMap(args []string) map[string]string {

	varMap := make(map[string]string)

	for _, fullArg := range args {
		splitArg := strings.SplitN(fullArg, "=", 2)
		if len(splitArg) < 2 {
			continue
		}

		varMap[splitArg[0]] = splitArg[1]
	}

	return varMap
}
