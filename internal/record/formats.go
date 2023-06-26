package record

// Ways to format the data in a set of records.

import (
	"fmt"
	"strings"
)

// Format a set of records as an env file.
//
//	export $(labrador --quiet | xargs)
func RecordsAsEnvFile(records map[string]*Record) (string, error) {

	result := ""

	for name, item := range records {
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
