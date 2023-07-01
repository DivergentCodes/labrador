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
	"github.com/divergentcodes/labrador/internal/record"
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

	// aws-ssmps
	defaultAwsSsmParameters := viper.GetViper().GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	fetchCmd.PersistentFlags().StringSlice("aws-param", defaultAwsSsmParameters, "AWS SSM parameter store path prefix or ARN")
	err = viper.BindPFlag(core.OptStr_AWS_SsmParameterStore, fetchCmd.PersistentFlags().Lookup("aws-param"))
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(fetchCmd)
}

// Top level logic for the fetch CLI subcommand
func fetch(cmd *cobra.Command, args []string) {
	ShowBanner()

	if countRemoteTargets() == 0 {
		core.PrintFatal("no remote values to fetch were specified", 1)
	}

	records := make(map[string]*record.Record, 0)

	records = fetchAwsSsmParameters(records)
	core.PrintNormal(fmt.Sprintf("\nFetched %d values\n", len(records)))

	formattedOutput := formatRecordsOutput(records)

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

	return remoteTargetCount
}

// Fetch values, convert to records, add to list, and return the list.
func fetchAwsSsmParameters(records map[string]*record.Record) map[string]*record.Record {

	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	if len(awsSsmParameters) != 0 {
		ssmRecords, err := aws.FetchParameterStore()
		if err != nil {
			core.PrintFatal("failed to get SSM parameters", 1)
		}

		core.PrintVerbose(fmt.Sprintf("\nFetched %d SSM parameters", len(ssmRecords)))
		for name, record := range ssmRecords {
			records[name] = record
			core.PrintVerbose(fmt.Sprintf("\n\t%s", record.Data["arn"]))
			core.PrintDebug(fmt.Sprintf("\n\t\ttype: \t\t%s", record.Data["type"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tversion: \t%s", record.Data["version"]))
			core.PrintDebug(fmt.Sprintf("\n\t\tmodified: \t%s", record.Data["last-modified"]))
		}
	}

	return records
}

// Convert the list of records to formatted output.
func formatRecordsOutput(records map[string]*record.Record) string {
	// Only does env format for now.
	formattedOutput, err := record.RecordsAsEnvFile(records)
	if err != nil {
		core.PrintFatal("failed to format records as env file", 1)
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
