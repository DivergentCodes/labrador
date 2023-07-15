package variable

// Ways to format the data in a set of variables.

import (
	"fmt"
	"strings"
)

// Format a set of variables as an env file.
//
//	export $(labrador --quiet | xargs)
func VariablesAsEnvFile(variables map[string]*Variable) (string, error) {

	result := ""

	for name, item := range variables {
		envVarName := envNamify(name)
		result += fmt.Sprintf("%s=%s\n", envVarName, item.Value)
	}

	return result, nil
}

// Transform strings into valid environment variable names.
func envNamify(name string) string {
	envVarName := name
	for _, c := range [2]string{"-", " "} {
		envVarName = strings.Replace(envVarName, c, "_", -1)
	}
	return envVarName
}
