/*
Labrador fetches variables and secrets from remote services.

Values are recursively pulled from one or more services, and output to the
terminal or a file.

Labrador is focused on reading and pulling values, not on managing or writing
values. It was created with CI/CD pipelines and network services in mind.

Usage:

	labrador [command]

Available Commands:

	completion  Generate the autocompletion script for the specified shell
	export      Fetch and export values as shell environment variables
	fetch       Fetch values from services
	help        Help about any command
	version     Print the version

Flags:

	-c, --config string   config file (default is .labrador.yaml)
	    --debug           Enable debug mode
	-h, --help            help for labrador
	-q, --quiet           Quiet CLI output
	    --verbose         Verbose CLI output

Use "labrador [command] --help" for more information about a command.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/divergentcodes/labrador/internal/aws"
	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/variable"
)

var (
	rootCmd = &cobra.Command{
		Use:   "labrador",
		Short: "Fetch and load variables and secrets from remote services",
		Long: bannerText() + `
Labrador fetches variables and secrets from remote services.

Values are recursively pulled from one or more services, and output to the
terminal or a file.

Labrador is focused on reading and pulling values, not on managing or writing
values. It was created with CI/CD pipelines and network services in mind.`,
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

	// json
	// TODO: implement exporting results to JSON.
	/*
		defaultOutJSON := viper.GetBool(core.OptStr_OutJSON)
		rootCmd.PersistentFlags().Bool("json", defaultOutJSON, "Use JSON output")
		err = viper.BindPFlag(core.OptStr_OutJSON, rootCmd.PersistentFlags().Lookup("json"))
		if err != nil {
			panic(err)
		}
	*/

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

// Count the number of user-defined resources to pull values from.
func countRemoteTargets() int {
	remoteTargetCount := 0

	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	remoteTargetCount += len(awsSsmParameters)

	awsSmSecrets := viper.GetStringSlice(core.OptStr_AWS_SecretManager)
	remoteTargetCount += len(awsSmSecrets)

	return remoteTargetCount
}

// Fetch AWS SSM Parameter Store values, convert to variables, add to list, and return the list.
func fetchAwsSsmParameters(variables map[string]*variable.Variable) map[string]*variable.Variable {

	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	if len(awsSsmParameters) != 0 {
		ssmVariables, err := aws.FetchParameterStore()
		if err != nil {
			core.PrintFatal("failed to get SSM parameters", 1)
		}

		core.PrintVerbose(fmt.Sprintf("\nFetched %d values from AWS SSM Parameter Store", len(ssmVariables)))
		for name, variable := range ssmVariables {
			variables[name] = variable
			core.PrintVerbose(fmt.Sprintf("\n\t%s", variable.Metadata["arn"]))
			core.PrintDebug(fmt.Sprintf("\n\t\ttype: \t\t%s", variable.Metadata["type"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tversion: \t%s", variable.Metadata["version"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tmodified: \t%s", variable.Metadata["last-modified"]))
		}
	}

	return variables
}

// Fetch AWS Secrets Manager values, convert to variables, add to list, and return the list.
func fetchAwsSmSecrets(variables map[string]*variable.Variable) map[string]*variable.Variable {

	awsSmSecrets := viper.GetStringSlice(core.OptStr_AWS_SecretManager)
	if len(awsSmSecrets) != 0 {
		smVariables, err := aws.FetchSecretsManager()
		if err != nil {
			core.PrintFatal("failed to get Secrets Manager values", 1)
		}

		core.PrintVerbose(fmt.Sprintf("\nFetched %d values from AWS Secrets Manager", len(smVariables)))
		for name, variable := range smVariables {
			variables[name] = variable
			core.PrintVerbose(fmt.Sprintf("\n\t%s (%s)", variable.Metadata["arn"], variable.Key))
			core.PrintDebug(fmt.Sprintf("\n\t\tsecret-name: \t%s", variable.Metadata["secret-name"]))
			core.PrintDebug(fmt.Sprintf("\n\t\ttype: \t\t%s", variable.Metadata["type"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tversion-id: \t%s", variable.Metadata["version-id"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tcreated: \t%s", variable.Metadata["created-date"]))
		}
	}

	return variables
}
