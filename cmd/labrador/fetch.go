package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/divergentcodes/labrador/internal/aws"
	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/variable"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch values from services",
	Long:  `Fetch values from services`,
	Run:   fetch,
}

// Initialize the fetch CLI subcommand
func init() {

	// outfile
	defaultOutFile := viper.GetViper().GetString(core.OptStr_OutFile)
	fetchCmd.PersistentFlags().StringP("outfile", "o", defaultOutFile, "File path to write variable/value pairs to")
	err := viper.BindPFlag(core.OptStr_OutFile, fetchCmd.PersistentFlags().Lookup("outfile"))
	if err != nil {
		panic(err)
	}

	// outfile-mode
	defaultFileMode := viper.GetViper().GetString(core.OptStr_FileMode)
	fetchCmd.PersistentFlags().String("outfile-mode", defaultFileMode, "File permissions for newly created outfile")
	err = viper.BindPFlag(core.OptStr_FileMode, fetchCmd.PersistentFlags().Lookup("outfile-mode"))
	if err != nil {
		panic(err)
	}

	// export
	defaultExport := viper.GetBool(core.OptStr_Export)
	fetchCmd.PersistentFlags().Bool("export", defaultExport, "Format as sh environment variables to use with 'source'")
	err = viper.BindPFlag(core.OptStr_Export, fetchCmd.PersistentFlags().Lookup("export"))
	if err != nil {
		panic(err)
	}

	// quote
	defaultQuote := viper.GetBool(core.OptStr_Quote)
	fetchCmd.PersistentFlags().Bool("quote", defaultQuote, "Surround each value with doublequotes")
	err = viper.BindPFlag(core.OptStr_Quote, fetchCmd.PersistentFlags().Lookup("quote"))
	if err != nil {
		panic(err)
	}

	// aws-region
	defaultAwsRegion := ""
	fetchCmd.PersistentFlags().String("aws-region", defaultAwsRegion, "AWS region")
	err = viper.BindPFlag(core.OptStr_AWS_Region, fetchCmd.PersistentFlags().Lookup("aws-region"))
	if err != nil {
		panic(err)
	}

	// aws-param
	defaultAwsSsmParameters := viper.GetViper().GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	fetchCmd.PersistentFlags().StringSlice("aws-param", defaultAwsSsmParameters, "AWS SSM parameter store path prefix")
	err = viper.BindPFlag(core.OptStr_AWS_SsmParameterStore, fetchCmd.PersistentFlags().Lookup("aws-param"))
	if err != nil {
		panic(err)
	}

	// aws-secret
	defaultAwsSmSecrets := viper.GetViper().GetStringSlice(core.OptStr_AWS_SecretManager)
	fetchCmd.PersistentFlags().StringSlice("aws-secret", defaultAwsSmSecrets, "AWS Secrets Manager secret name")
	err = viper.BindPFlag(core.OptStr_AWS_SecretManager, fetchCmd.PersistentFlags().Lookup("aws-secret"))
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(fetchCmd)
}

// Top level logic for the fetch CLI subcommand
func fetch(cmd *cobra.Command, args []string) {
	// --export implies --quiet
	if viper.GetBool(core.OptStr_Export) {
		viper.Set(core.OptStr_Quiet, true)
	}

	ShowBanner()

	if countRemoteTargets() == 0 {
		core.PrintFatal("no remote values to fetch were specified", 1)
	}

	variables := make(map[string]*variable.Variable, 0)
	variables = fetchAwsSsmParameters(variables)
	variables = fetchAwsSmSecrets(variables)

	core.PrintDebug("\n")
	core.PrintNormal(fmt.Sprintf("\nFetched %d values\n", len(variables)))
	core.PrintDebug("\n")

	formattedOutput := formatVariablesOutput(variables)

	outFilePath := viper.GetString(core.OptStr_OutFile)
	outFileMode := viper.GetString(core.OptStr_FileMode)
	if outFilePath != "" {
		// Dump formatted results to file.
		writeFormattedOutFile(formattedOutput, outFilePath, outFileMode)
		core.PrintNormal(fmt.Sprintf("Wrote parameters to file: %s\n", outFilePath))
	} else {
		// Display formatted results to STDOUT.
		core.PrintAlways(formattedOutput)
		core.PrintNormal("\n")
	}
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

// Convert the list of variables to formatted output.
func formatVariablesOutput(variables map[string]*variable.Variable) string {
	var formattedOutput string
	var err error

	asExport := viper.GetBool(core.OptStr_Export)

	if asExport {
		formattedOutput, err = variable.VariablesAsShellExport(variables)
		if err != nil {
			core.PrintFatal("failed to format variables as shell exports", 1)
		}
	} else {
		useQuotes := viper.GetBool(core.OptStr_Quote)
		formattedOutput, err = variable.VariablesAsEnvFile(variables, useQuotes)
		if err != nil {
			core.PrintFatal("failed to format variables as env file", 1)
		}
	}

	return formattedOutput
}

// Write fetched, formatted values to file.
func writeFormattedOutFile(formattedOutput string, outFilePath string, outFileMode string) {
	outFilePath = filepath.Clean(outFilePath)

	modeValue, _ := strconv.ParseUint(outFileMode, 8, 32)
	fileMode := os.FileMode(modeValue)

	// If the file doesn't exist, create it, or append to the file.
	fh, err := os.OpenFile(outFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode) //#nosec
	if err != nil {
		core.PrintFatal(err.Error(), 1)
	}

	defer fh.Close()

	if _, err = fh.WriteString(formattedOutput); err != nil {
		core.PrintFatal(err.Error(), 1)
	}
}
