// Package gcp fetches values from GCP Secret Manager.
package gcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/grpc/status"

	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/variable"
	"github.com/spf13/viper"
)

// Fetch values from GCP Secret Manager.
func FetchSecretManager() (map[string]*variable.Variable, error) {

	ctx := context.Background()
	smClient := initSecretManagerClient(&ctx)
	defer smClient.Close()

	secretManagerResources := viper.GetStringSlice(core.OptStr_GCP_SecretManager)
	secretManagerVariables := make(map[string]*variable.Variable, 0)

	core.PrintVerbose("\nFetching GCP Secret Manager values...")
	for _, resource := range secretManagerResources {
		core.PrintDebug(fmt.Sprintf("\n\t%s", resource))
	}

	// Fetch and aggregate the resources.
	for _, resource := range secretManagerResources {

		// Handle resource string with and without explicit version.
		pathParts := strings.Split(resource, "/")
		secretName := resource
		secretVersion := "versions/latest"
		if len(pathParts) == 6 {
			secretName = strings.Join(pathParts[0:4], "/")
			secretVersion = strings.Join(pathParts[4:], "/")
		}
		secretFullPath := fmt.Sprintf("%s/%s", secretName, secretVersion)

		// Fetch data from GCP.
		secretValue := getSecretValue(&ctx, smClient, secretFullPath)
		secretMetadata := getSecretMetadata(&ctx, smClient, secretName)
		secretManagerVariables = secretToVariable(secretValue, secretMetadata, secretManagerVariables)
	}

	return secretManagerVariables, nil
}

// Initialize a GCP Secret Manager client instance.
func initSecretManagerClient(ctx *context.Context) *secretmanager.Client {

	core.PrintDebug("\n")
	core.PrintVerbose("\nInitializing GCP Secret Manager client...")

	// Read service account credentials from JSON file pointed at by $GOOGLE_APPLICATION_CREDENTIALS.
	smClient, err := secretmanager.NewClient(*ctx)
	if err != nil {
		log.Fatalf("failed to initialize GCP Secret Manager client, %v", err)
		if s, ok := status.FromError(err); ok {
			log.Println(s.Message())
			for _, d := range s.Proto().Details {
				log.Println(d)
			}
		}
	}

	return smClient
}

// Get the value of the GCP Secret Manager secret.
func getSecretValue(ctx *context.Context, smClient *secretmanager.Client, secretFullPath string) *secretmanagerpb.AccessSecretVersionResponse {
	request := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretFullPath,
	}

	// Call the API.
	result, err := smClient.AccessSecretVersion(*ctx, request)
	if err != nil {
		log.Fatalf("\nfailed to fetch GCP Secret Manager value, %v", err)
		if s, ok := status.FromError(err); ok {
			log.Println(s.Message())
			for _, d := range s.Proto().Details {
				log.Println(d)
			}
		}
	}

	return result
}

// Get the metadata of the GCP Secret Manager secret.
func getSecretMetadata(ctx *context.Context, smClient *secretmanager.Client, secretName string) *secretmanagerpb.Secret {
	request := &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	}
	result, err := smClient.GetSecret(*ctx, request)
	if err != nil {
		log.Fatalf("\nfailed to fetch GCP Secret Manager metadata, %v", err)
		if s, ok := status.FromError(err); ok {
			log.Println(s.Message())
			for _, d := range s.Proto().Details {
				log.Println(d)
			}
		}
	}
	return result
}

// Convert a GCP Secret Manager secret to a Variable.
func secretToVariable(secretValue *secretmanagerpb.AccessSecretVersionResponse, secretMetadata *secretmanagerpb.Secret, smSecretVariables map[string]*variable.Variable) map[string]*variable.Variable {
	pathParts := strings.Split(secretValue.GetName(), "/")

	key := pathParts[3]
	value := string(secretValue.GetPayload().GetData())

	project := pathParts[1]
	version := pathParts[5]
	secretName := strings.Join(pathParts[0:4], "/")

	var annotationSlice []string
	for k, v := range secretMetadata.Annotations {
		annotationSlice = append(annotationSlice, fmt.Sprintf("%s=%s", k, v))
	}
	annotationString := strings.Join(annotationSlice, ",")

	var labelSlice []string
	for k, v := range secretMetadata.Labels {
		labelSlice = append(labelSlice, fmt.Sprintf("%s=%s", k, v))
	}
	labelString := strings.Join(labelSlice, ",")

	var topicSlice []string
	for _, topic := range secretMetadata.GetTopics() {
		topicSlice = append(topicSlice, topic.Name)
	}
	topicString := strings.Join(topicSlice, ",")

	result := variable.Variable{
		Source:   "gcp-secret-manager",
		Key:      key,
		Value:    value,
		Metadata: make(map[string]string),
	}
	result.Metadata["secret-name"] = secretName
	result.Metadata["project"] = project
	result.Metadata["create-time"] = secretMetadata.GetCreateTime().AsTime().String() //CreateTime.AsTime().String()
	result.Metadata["expire-time"] = secretMetadata.GetExpireTime().AsTime().String()
	result.Metadata["version"] = version
	result.Metadata["etag"] = secretMetadata.Etag
	result.Metadata["rotation"] = secretMetadata.GetRotation().String()
	result.Metadata["topics"] = topicString
	result.Metadata["annotations"] = annotationString
	result.Metadata["labels"] = labelString

	smSecretVariables[key] = &result

	return smSecretVariables
}
