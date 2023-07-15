// Package aws fetches values from AWS SSM Parameter Store.
package aws

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/spf13/viper"

	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/variable"
)

// Fetch values from AWS SSM Parameter Store.
func FetchParameterStore() (map[string]*variable.Variable, error) {

	ssmClient := initSsmClient()
	ssmParameterResources := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	ssmParameterVariables := make(map[string]*variable.Variable, 0)

	core.PrintVerbose("\nFetching SSM Parameter Store values...")
	for _, resource := range ssmParameterResources {
		core.PrintDebug(fmt.Sprintf("\n\t%s", resource))
	}

	// Fetch and aggregate the parameter resources.
	for _, resource := range ssmParameterResources {
		if strings.HasSuffix(resource, "/*") {
			// Wildcard parameter paths.
			ssmParameterResultBatch := fetchParameterStoreWildcard(ssmClient, resource)
			for name, variable := range ssmParameterResultBatch {
				ssmParameterVariables[name] = variable
			}
		} else {
			// Single parameter paths.
			ssmParameterResultBatch := fetchParameterStoreSingle(ssmClient, resource)
			for name, variable := range ssmParameterResultBatch {
				ssmParameterVariables[name] = variable
			}
		}
	}

	return ssmParameterVariables, nil
}

// Initialize a SSM client.
func initSsmClient() *ssm.Client {
	awsRegion := viper.GetString(core.OptStr_AWS_Region)

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}
	if awsRegion != "" {
		awsConfig.Region = awsRegion
		core.PrintDebug(fmt.Sprintf("\nSet AWS region: %s", awsRegion))
	}

	core.PrintVerbose("\nInitializing AWS SSM client...")
	ssmClient := ssm.NewFromConfig(awsConfig)
	if err != nil {
		log.Fatalf("failed to initialize AWS SSM client, %v", err)
	}

	return ssmClient
}

// Recursively fetch all parameters at a SSM parameter store wildcard path.
func fetchParameterStoreSingle(ssmClient *ssm.Client, resource string) map[string]*variable.Variable {

	// Using a map to be consistent with the wilcard fetching.
	ssmParameterResults := make(map[string]*variable.Variable, 0)

	input := &ssm.GetParameterInput{
		Name:           aws.String(resource),
		WithDecryption: aws.Bool(true),
	}

	resp, err := ssmClient.GetParameter(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to fetch AWS SSM Parameter Store values, %v", err)
	}

	// Convert the result to a canonical variable.
	result := parameterToVariable(resp.Parameter)
	ssmParameterResults[result.Key] = result

	return ssmParameterResults
}

// Recursively fetch all parameters at a SSM parameter store wildcard path.
func fetchParameterStoreWildcard(ssmClient *ssm.Client, resource string) map[string]*variable.Variable {

	recursive := true
	nextToken := ""
	ssmParameterResults := make(map[string]*variable.Variable, 0)

	resource = strings.TrimRight(resource, "/*")

	// Only 10 parameters can be fetched per call. Loop to fetch all.
	for {
		input := &ssm.GetParametersByPathInput{
			Path:           aws.String(resource),
			Recursive:      aws.Bool(recursive),
			WithDecryption: aws.Bool(true),
			MaxResults:     aws.Int32(10),
			NextToken:      aws.String(nextToken),
		}

		// Fetch the parameters.
		resp, err := ssmClient.GetParametersByPath(context.TODO(), input)
		if err != nil {
			log.Fatalf("failed to fetch SSM parameters, %v", err)
		}

		// Aggregate the parameters, since the call can be recursive.
		// Last variable has highest precendence.
		for i := range resp.Parameters {
			result := parameterToVariable(&resp.Parameters[i])
			ssmParameterResults[result.Key] = result
		}

		// Determine if all parameters have been fetched.
		if resp.NextToken == nil {
			break
		}
		nextToken = *resp.NextToken
	}

	return ssmParameterResults
}

// Convert a parameter store resource to an intermediate labrador variable representation.
func parameterToVariable(parameter *ssmTypes.Parameter) *variable.Variable {

	splitArn := strings.Split(*parameter.ARN, "/")
	varKey := splitArn[len(splitArn)-1]

	result := variable.Variable{
		Source:   "aws-ssm-parameter-store",
		Key:      varKey,
		Value:    *parameter.Value,
		Metadata: make(map[string]string),
	}
	result.Metadata["arn"] = *parameter.ARN
	result.Metadata["path"] = *parameter.Name
	result.Metadata["type"] = string(parameter.Type)
	result.Metadata["last-modified"] = parameter.LastModifiedDate.String()
	result.Metadata["version"] = fmt.Sprintf("%d", parameter.Version)

	return &result
}
