package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"

	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/record"
)

// Fetch values from AWS Secrets Manager.
func FetchSecretsManager() (map[string]*record.Record, error) {

	smClient := initSecretsManagerClient()
	secretsManagerResources := viper.GetStringSlice(core.OptStr_AWS_SecretManager)
	secretsManagerRecords := make(map[string]*record.Record, 0)

	core.PrintVerbose("\nFetching Secrets Manager values...")
	for _, resource := range secretsManagerResources {
		core.PrintDebug(fmt.Sprintf("\n\t%s", resource))
	}

	// Fetch and aggregate the parameter resources.
	for _, resource := range secretsManagerResources {
		smSecretsManagerResultBatch := fetchSecretsManagerSecret(smClient, resource)
		for name, record := range smSecretsManagerResultBatch {
			secretsManagerRecords[name] = record
		}
	}

	return secretsManagerRecords, nil
}

// Initialize a AWS Secrets Manager client instance.
func initSecretsManagerClient() *secretsmanager.Client {
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

	core.PrintDebug("\n")
	core.PrintVerbose("\nInitializing AWS Secrets Manager client...")
	smClient := secretsmanager.NewFromConfig(awsConfig)
	if err != nil {
		log.Fatalf("failed to initialize AWS Secrets Manager client, %v", err)
	}

	return smClient
}

// Fetch a secret from AWS Secrets Manager.
func fetchSecretsManagerSecret(smClient *secretsmanager.Client, resource string) map[string]*record.Record {
	// Using a map to be consistent with the wilcard fetching.
	smSecretResults := make(map[string]*record.Record, 0)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(resource),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	resp, err := smClient.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to fetch AWS Secrets Manager values, %v", err)
	}

	smSecretResults = secretToRecords(resp, smSecretResults)

	return smSecretResults
}

// Convert an AWS Secrets Manager secret to a list of Records.
//
// One secret can hold multiple key/value pairs.
func secretToRecords(secret *secretsmanager.GetSecretValueOutput, smSecretRecords map[string]*record.Record) map[string]*record.Record {

	var varType string
	if secret.SecretString != nil {
		varType = "SecretString"

		// Extract key/value pairs from JSON.
		var secretDict map[string]string
		err := json.Unmarshal([]byte(*secret.SecretString), &secretDict)
		if err != nil {
			core.PrintFatal(err.Error(), 1)
		}

		// Format each key/value pair as a record.
		for k, v := range secretDict {
			result := record.Record{
				Source:   "aws-secrets-manager",
				Key:      k,
				Value:    v,
				Metadata: make(map[string]string),
			}
			result.Metadata["arn"] = *secret.ARN
			result.Metadata["secret-name"] = *secret.Name
			result.Metadata["type"] = varType
			result.Metadata["created-date"] = secret.CreatedDate.String()
			result.Metadata["version-id"] = *secret.VersionId
			//result.Metadata["version-stages"] = *&secret.VersionStages[]

			smSecretRecords[k] = &result
		}
	} else {
		varType = "SecretBinary"

		result := record.Record{
			Source:   "aws-secrets-manager",
			Key:      *secret.Name,
			Value:    string(secret.SecretBinary[:]),
			Metadata: make(map[string]string),
		}
		result.Metadata["arn"] = *secret.ARN
		result.Metadata["secret-name"] = *secret.Name
		result.Metadata["type"] = varType
		result.Metadata["created-date"] = secret.CreatedDate.String()
		result.Metadata["version-id"] = *secret.VersionId
		//result.Metadata["version-stages"] = *&secret.VersionStages[]

		smSecretRecords[*secret.Name] = &result
	}

	return smSecretRecords
}
