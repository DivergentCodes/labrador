package variable

// Ways to format the data in a set of variables.

import (
	"fmt"
	"strings"
)

// Format a set of variables as an env file.
//
//	export $(labrador --quiet | xargs)
func VariablesAsEnvFile(variables map[string]*Variable, quote bool) (string, error) {

	result := ""

	for name, item := range variables {
		envVarName := envNamify(name)
		envVarValue := item.Value
		if quote {
			envVarValue = escapeDoubleQuotes(envVarValue)
			envVarValue = fmt.Sprintf("\"%s\"\n", envVarValue)
		}
		result += fmt.Sprintf("%s=%s\n", envVarName, envVarValue)
	}

	return result, nil
}

// Escape double quotes in provided string.
func escapeDoubleQuotes(value string) string {
	envVarValue := strings.Replace(value, "\"", "\\\"", -1)
	return envVarValue
}

// Transform strings into valid environment variable names.
func envNamify(name string) string {
	envVarName := name
	for _, c := range [2]string{"-", " "} {
		envVarName = strings.Replace(envVarName, c, "_", -1)
	}
	return envVarName
}
