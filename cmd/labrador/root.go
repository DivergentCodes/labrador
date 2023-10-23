/*
Labrador (latest) created by Jack Sullivan <jack@divergent.codes>

Labrador is a CLI tool to fetch secrets and other configuration values
from one or more remote services.

Labrador was created to explore safer, consistent, cross-platform ways of
handling secrets during each phase of the SDLC. The idea is to fetch secrets
from a central service at runtime, in a standard way, instead of copying secrets
to each environment and persisting them all over the place.

Usage:

	labrador [command]

Available Commands:

	completion  Generate the autocompletion script for the specified shell
	export      Fetch and export values as shell environment variables
	fetch       Fetch values from services
	help        Help about any command
	version     Print the version

Flags:

	    --aws-param strings    AWS SSM parameter store path prefix
	    --aws-region string    AWS region
	    --aws-secret strings   AWS Secrets Manager secret name
	-c, --config string        config file (default is .labrador.yaml)
	    --debug                Enable debug mode
	-h, --help                 help for labrador
	    --lower                Set all variable names to lower case
	-q, --quiet                Quiet CLI output
	    --quote                Surround each value with doublequotes
	    --upper                Set all variable names to upper case
	    --verbose              Verbose CLI output

Use "labrador [command] --help" for more information about a command.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/divergentcodes/labrador/internal/aws"
	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/gcp"
	"github.com/divergentcodes/labrador/internal/variable"
)

var (
	rootCmd = &cobra.Command{
		Use:   "labrador",
		Short: "Fetch and secrets and values from remote services",
		Long: bannerText() + `
Labrador is a CLI tool to fetch secrets and other configuration values
from one or more remote services.

Labrador was created to explore safer, consistent, cross-platform ways of
handling secrets during each phase of the SDLC. The idea is to fetch secrets
from a central service at runtime, in a standard way, instead of copying secrets
to each environment and persisting them all over the place.
`,
	}
)

// Run for root command. Will execute for all subcommands.
func init() {
	cobra.OnInitialize(core.InitConfigInstance)
	initRootFlags()
}

func initRootFlags() {

	// config
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is .labrador.yaml)")
	err := viper.BindPFlag(core.OptStr_Config, rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		panic(err)
	}

	// Display.

	// debug
	defaultDebug := viper.GetBool(core.OptStr_Debug)
	rootCmd.PersistentFlags().Bool("debug", defaultDebug, "Enable debug mode")
	err = viper.BindPFlag(core.OptStr_Debug, rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		panic(err)
	}

	// quiet
	defaultQuiet := viper.GetBool(core.OptStr_Quiet)
	rootCmd.PersistentFlags().BoolP("quiet", "q", defaultQuiet, "Quiet CLI output")
	err = viper.BindPFlag(core.OptStr_Quiet, rootCmd.PersistentFlags().Lookup("quiet"))
	if err != nil {
		panic(err)
	}

	// verbose
	defaultVerbose := viper.GetBool(core.OptStr_Verbose)
	rootCmd.PersistentFlags().Bool("verbose", defaultVerbose, "Verbose CLI output")
	err = viper.BindPFlag(core.OptStr_Verbose, rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		panic(err)
	}

	// Output transforms.

	// quote
	defaultQuote := viper.GetBool(core.OptStr_Quote)
	rootCmd.PersistentFlags().Bool("quote", defaultQuote, "Surround each value with doublequotes")
	err = viper.BindPFlag(core.OptStr_Quote, rootCmd.PersistentFlags().Lookup("quote"))
	if err != nil {
		panic(err)
	}

	// lower
	defaultToLower := viper.GetBool(core.OptStr_ToLower)
	rootCmd.PersistentFlags().Bool("lower", defaultToLower, "Set all variable names to lower case")
	err = viper.BindPFlag(core.OptStr_ToLower, rootCmd.PersistentFlags().Lookup("lower"))
	if err != nil {
		panic(err)
	}

	// upper
	defaultToUpper := viper.GetBool(core.OptStr_ToUpper)
	rootCmd.PersistentFlags().Bool("upper", defaultToUpper, "Set all variable names to upper case")
	err = viper.BindPFlag(core.OptStr_ToUpper, rootCmd.PersistentFlags().Lookup("upper"))
	if err != nil {
		panic(err)
	}

	// Remote services.

	// aws-region
	defaultAwsRegion := ""
	rootCmd.PersistentFlags().String("aws-region", defaultAwsRegion, "AWS region")
	err = viper.BindPFlag(core.OptStr_AWS_Region, rootCmd.PersistentFlags().Lookup("aws-region"))
	if err != nil {
		panic(err)
	}

	// aws-param
	defaultAwsSsmParameters := viper.GetViper().GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	rootCmd.PersistentFlags().StringSlice("aws-param", defaultAwsSsmParameters, "AWS SSM parameter store path prefix")
	err = viper.BindPFlag(core.OptStr_AWS_SsmParameterStore, rootCmd.PersistentFlags().Lookup("aws-param"))
	if err != nil {
		panic(err)
	}

	// aws-secret
	defaultAwsSmSecrets := viper.GetViper().GetStringSlice(core.OptStr_AWS_SecretsManager)
	rootCmd.PersistentFlags().StringSlice("aws-secret", defaultAwsSmSecrets, "AWS Secrets Manager secret name")
	err = viper.BindPFlag(core.OptStr_AWS_SecretsManager, rootCmd.PersistentFlags().Lookup("aws-secret"))
	if err != nil {
		panic(err)
	}

	// gcp-secret
	defaultGcpSecret := viper.GetViper().GetStringSlice(core.OptStr_GCP_SecretManager)
	rootCmd.PersistentFlags().StringSlice("gcp-secret", defaultGcpSecret, "GCP Secret Manager secret name")
	err = viper.BindPFlag(core.OptStr_GCP_SecretManager, rootCmd.PersistentFlags().Lookup("gcp-secret"))
	if err != nil {
		panic(err)
	}

	rootCmd.MarkFlagsMutuallyExclusive("lower", "upper")
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "debug")
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "verbose")
}

// Execute starts the app CLI.
func Execute() error {
	return rootCmd.Execute()
}

// ShowBanner shows the CLI banner message.
func ShowBanner() {

	core.PrintNormal(bannerText())

	core.PrintDebug("Labrador DEBUG mode is enabled\n")

	msg := fmt.Sprintf(
		"\tVersion: %s\n\tCommit: %s\n\tBuilt on: %s\n\tBuilt by: %s\n",
		core.Version, core.Commit, core.Date, core.BuiltBy,
	)
	core.PrintDebug(msg)
}

func bannerText() string {
	return fmt.Sprintf("Labrador %s created by %s <%s>\n", core.Version, core.AuthorName, core.AuthorEmail)
}

// Wrapper to bundle fetching from each service.
func fetchVariables() map[string]*variable.Variable {
	if countRemoteTargets() == 0 {
		core.PrintFatal("no remote values to fetch were specified", 1)
	}

	variables := make(map[string]*variable.Variable, 0)

	variables = fetchAwsSsmParameters(variables)
	variables = fetchAwsSmSecrets(variables)

	variables = fetchGcpSmSecrets(variables)

	return variables
}

// Count the number of user-defined resources to pull values from.
func countRemoteTargets() int {
	remoteTargetCount := 0

	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	remoteTargetCount += len(awsSsmParameters)

	awsSmSecrets := viper.GetStringSlice(core.OptStr_AWS_SecretsManager)
	remoteTargetCount += len(awsSmSecrets)

	return remoteTargetCount
}

// Fetch AWS SSM Parameter Store values, convert to variables, add to list, and return the list.
func fetchAwsSsmParameters(variables map[string]*variable.Variable) map[string]*variable.Variable {

	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	if len(awsSsmParameters) != 0 {
		ssmVariables, err := aws.FetchParameterStore()
		if err != nil {
			core.PrintFatal("failed to get AWS SSM parameters", 1)
		}

		core.PrintVerbose(fmt.Sprintf("\nFetched %d values from AWS SSM Parameter Store", len(ssmVariables)))
		for name, variable := range ssmVariables {
			variables[name] = variable
			core.PrintVerbose(fmt.Sprintf("\n\t%s (%s)", variable.Metadata["arn"], variable.Key))
			core.PrintDebug(fmt.Sprintf("\n\t\tarn: \t\t%s", variable.Metadata["arn"]))
			core.PrintDebug(fmt.Sprintf("\n\t\ttype: \t\t%s", variable.Metadata["type"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tversion: \t%s", variable.Metadata["version"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tmodified: \t%s", variable.Metadata["last-modified"]))
		}
	}

	return variables
}

// Fetch AWS Secrets Manager values, convert to variables, add to list, and return the list.
func fetchAwsSmSecrets(variables map[string]*variable.Variable) map[string]*variable.Variable {

	awsSmSecrets := viper.GetStringSlice(core.OptStr_AWS_SecretsManager)
	if len(awsSmSecrets) != 0 {
		smVariables, err := aws.FetchSecretsManager()
		if err != nil {
			core.PrintFatal("failed to get AWS Secrets Manager values", 1)
		}

		core.PrintVerbose(fmt.Sprintf("\nFetched %d values from AWS Secrets Manager", len(smVariables)))
		for name, variable := range smVariables {
			variables[name] = variable
			core.PrintVerbose(fmt.Sprintf("\n\t%s (%s)", variable.Metadata["arn"], variable.Key))
			core.PrintDebug(fmt.Sprintf("\n\t\tarn: \t\t%s", variable.Metadata["arn"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tsecret-name: \t%s", variable.Metadata["secret-name"]))
			core.PrintDebug(fmt.Sprintf("\n\t\ttype: \t\t%s", variable.Metadata["type"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tversion-id: \t%s", variable.Metadata["version-id"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tcreated: \t%s", variable.Metadata["created-date"]))
		}
	}

	return variables
}

// Fetch GCP Secret Manager values, convert to variables, add to list, and return the list.
func fetchGcpSmSecrets(variables map[string]*variable.Variable) map[string]*variable.Variable {

	gcpSmSecrets := viper.GetStringSlice(core.OptStr_GCP_SecretManager)
	if len(gcpSmSecrets) != 0 {
		smVariables, err := gcp.FetchSecretManager()
		if err != nil {
			core.PrintFatal("failed to get GCP Secret Manager values", 1)
		}

		core.PrintVerbose(fmt.Sprintf("\nFetched %d values from GCP Secret Manager", len(smVariables)))
		for name, variable := range smVariables {
			variables[name] = variable
			core.PrintVerbose(fmt.Sprintf("\n\t%s (%s)", variable.Metadata["secret-name"], variable.Key))
			core.PrintDebug(fmt.Sprintf("\n\t\tsecret-name: \t%s", variable.Metadata["secret-name"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tcreate-time: \t%s", variable.Metadata["create-time"]))
			core.PrintDebug(fmt.Sprintf("\n\t\texpire-time: \t%s", variable.Metadata["expire-time"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tversion: \t%s", variable.Metadata["version"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tproject: \t%s", variable.Metadata["project"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tetag: \t\t%s", variable.Metadata["etag"]))
			core.PrintDebug(fmt.Sprintf("\n\t\trotation: \t%s", variable.Metadata["rotation"]))
			core.PrintDebug(fmt.Sprintf("\n\t\ttopics: \t%s", variable.Metadata["topics"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tannotations: \t%s", variable.Metadata["annotations"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tlabels: \t%s", variable.Metadata["labels"]))
		}
	}

	return variables
}
