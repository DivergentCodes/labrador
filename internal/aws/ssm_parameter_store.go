// Package aws fetches values from AWS SSM Parameter Store.
package aws

import (
	"context"
	"fmt"
	"log"
	"os"
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

	core.PrintVerbose("\nFetching SSM parameters:")
	for _, resource := range ssmParameterResources {
		core.PrintVerbose(fmt.Sprintf("\n\t%s", resource))
	}

	// Fetch and aggregate the parameter resources.
	for _, resource := range ssmParameterResources {
		ssmParameterResultBatch := fetchParameterStoreResource(ssmClient, resource)
		for name, record := range ssmParameterResultBatch {
			ssmParameterRecords[name] = record
		}
	}

	return ssmParameterRecords, nil
}

// Initialize a SSM client.
func initSsmClient() *ssm.Client {
	//verbose := viper.GetBool(core.OptStr_Verbose)

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

// Given a SSM parameter path or ARN, fetch and return the parameter.
func fetchParameterStoreResource(ssmClient *ssm.Client, resource string) map[string]*record.Record {
	recursive := true
	nextToken := ""
	ssmParameterResults := make(map[string]*record.Record, 0)

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

// This isn't used right now. Scraps to use for later.
func SsmToFile() {

	OUTFILE := "test_out.txt"

	m := make(map[string]string)

	// Write parameters to file.
	fh, err := os.OpenFile(OUTFILE, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer fh.Close()

	// Isolate variable name from full parameter store path.
	// Then write it to file.
	// TODO: separate the variable name extraction from file writing.
	for full_key, value := range m {

		split_key := strings.Split(full_key, "/")
		var_name := split_key[len(split_key)-1]

		line := fmt.Sprintf("%s=%s\n", var_name, value)
		if _, err = fh.WriteString(line); err != nil {
			panic(err)
		}

	}
	core.PrintNormal(fmt.Sprintf("\nWrote parameters to file: %s\n", OUTFILE))
}
