package variable

// Ways to format the data in a set of variables.

import (
	"fmt"
	"regexp"
	"strings"
)

// Format a set of variables as an env file.
func VariablesAsEnvFile(variables map[string]*Variable, quote bool, lower bool, upper bool) (string, error) {

	result := ""

	for name, item := range variables {
		envVarName := envNamify(name)
		if lower {
			envVarName = strings.ToLower(envVarName)
		} else if upper {
			envVarName = strings.ToUpper(envVarName)
		}

		envVarValue := item.Value
		if quote {
			envVarValue = escapeDoubleQuotes(envVarValue)
			envVarValue = fmt.Sprintf("\"%s\"", envVarValue)
		}

		result += fmt.Sprintf("%s=%s\n", envVarName, envVarValue)
	}
	result = strings.TrimSuffix(result, "\n")

	return result, nil
}

// Format a set of shell environment variable exports.
//
// source <(labrador export)
func VariablesAsShellExport(variables map[string]*Variable, lower bool, upper bool) (string, error) {

	result, err := VariablesAsEnvFile(variables, true, lower, upper)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?m)^`)
	result = re.ReplaceAllString(result, "export ")

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
