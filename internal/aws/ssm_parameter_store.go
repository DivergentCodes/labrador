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
	"github.com/divergentcodes/labrador/internal/record"
)

func FetchParameterStore() (map[string]*record.Record, error) {

	ssmClient := initSsmClient()
	ssmParameterResources := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	ssmParameterRecords := make(map[string]*record.Record, 0)

	core.PrintVerbose("\nFetching SSM parameters...")
	for _, resource := range ssmParameterResources {
		core.PrintDebug(fmt.Sprintf("\n\t%s", resource))
	}

	// Fetch and aggregate the parameter resources.
	for _, resource := range ssmParameterResources {
		if strings.HasSuffix(resource, "/*") {
			// Wildcard parameter paths.
			ssmParameterResultBatch := fetchParameterStoreWildcard(ssmClient, resource)
			for name, record := range ssmParameterResultBatch {
				ssmParameterRecords[name] = record
			}
		} else {
			// Single parameter paths.
			ssmParameterResultBatch := fetchParameterStoreSingle(ssmClient, resource)
			for name, record := range ssmParameterResultBatch {
				ssmParameterRecords[name] = record
			}
		}
	}

	return ssmParameterRecords, nil
}

// Initialize a SSM client.
func initSsmClient() *ssm.Client {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	core.PrintVerbose("\nInitializing SSM client...")
	ssmClient := ssm.NewFromConfig(awsConfig)
	if err != nil {
		log.Fatalf("failed to initialize SSM client, %v", err)
	}

	return ssmClient
}

// Recursively fetch all parameters at a SSM parameter store wildcard path.
func fetchParameterStoreSingle(ssmClient *ssm.Client, resource string) map[string]*record.Record {

	// Using a map to be consistent with the wilcard fetching.
	ssmParameterResults := make(map[string]*record.Record, 0)

	input := &ssm.GetParameterInput{
		Name:           aws.String(resource),
		WithDecryption: aws.Bool(true),
	}

	resp, err := ssmClient.GetParameter(context.TODO(), input)

	if err != nil {
		log.Fatalf("failed to fetch SSM parameters, %v", err)
	}

	// Aggregate the parameters, since the call can be recursive.
	// Last record has highest precendence.
	result := parameterToRecord(resp.Parameter)
	ssmParameterResults[result.Key] = result

	return ssmParameterResults
}

// Recursively fetch all parameters at a SSM parameter store wildcard path.
func fetchParameterStoreWildcard(ssmClient *ssm.Client, resource string) map[string]*record.Record {

	recursive := true
	nextToken := ""
	ssmParameterResults := make(map[string]*record.Record, 0)

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
		// Last record has highest precendence.
		for i := range resp.Parameters {
			result := parameterToRecord(&resp.Parameters[i])
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

// Convert a parameter store resource to an intermediate labrador record representation.
func parameterToRecord(parameter *ssmTypes.Parameter) *record.Record {

	splitArn := strings.Split(*parameter.ARN, "/")
	paramKey := splitArn[len(splitArn)-1]

	result := record.Record{
		Source: "aws-ssm-parameter-store",
		Key:    paramKey,
		Value:  *parameter.Value,
		Data:   make(map[string]string),
	}
	result.Data["arn"] = *parameter.ARN
	result.Data["path"] = *parameter.Name
	result.Data["type"] = string(parameter.Type)
	result.Data["last-modified"] = parameter.LastModifiedDate.String()
	result.Data["version"] = fmt.Sprintf("%d", parameter.Version)

	return &result
}
