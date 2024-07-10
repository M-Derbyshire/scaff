package variable

import (
	"regexp"
	"strings"
)

// Populate returns the given string with the variable tags replaced with values from the given map.
// If a variable doesn't exist in the map, the user is prompted to provide it (and it is then added to the map)
// Once done, this will replace any escaped opening braces
func Populate(text string, vars map[string]string) (string, error) {
	resolvedText := text

	// Regex explantion:
	// Matches a series of alphanumeric characters surrounded by "{:" and ":}". The alphanumeric characters
	// can also be preceeded/proceeded by spaces.
	// Tags can be escaped by placing a backslash between the opening handlebar-brace and the colon ("{\:")
	varTagRegex := regexp.MustCompile(`{: *[a-zA-Z0-9-_]+ *:}`)

	for {
		// Get the first variable tag
		variableTag := varTagRegex.FindString(resolvedText)
		if variableTag == "" {
			break // No more tags
		}

		// Get the variable name out of the tag
		variableName := strings.Replace(variableTag, "{:", "", 1)
		variableName = strings.Replace(variableName, ":}", "", 1)
		variableName = strings.TrimSpace(variableName)

		// Resolve the variable value
		variableValue, varExists := vars[variableName]
		if !varExists {
			newVariableValue, err := Prompt(variableName)
			if err != nil {
				return "", err
			}

			variableValue = newVariableValue
			vars[variableName] = variableValue
		}

		// Now replace the tag
		resolvedText = strings.Replace(resolvedText, variableTag, variableValue, 1)
	}

	// Finally, resolve any escaped tags
	resolvedText = strings.ReplaceAll(resolvedText, "{\\:", "{:")

	return resolvedText, nil
}
