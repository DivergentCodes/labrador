package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"divergent.codes/labrador/internal/aws"
	"divergent.codes/labrador/internal/core"
	"divergent.codes/labrador/internal/record"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch values from services",
	Long:  `Fetch values from services`,
	Run:   fetch,
}

// Initialize the fetch CLI subcommand
func init() {

	defaultAwsSsmParameters := viper.GetViper().GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	fetchCmd.PersistentFlags().StringSlice("aws-ps", defaultAwsSsmParameters, "AWS SSM parameter store path prefix or ARN")
	err := viper.BindPFlag(core.OptStr_AWS_SsmParameterStore, fetchCmd.PersistentFlags().Lookup(core.OptStr_AWS_SsmParameterStore))
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(fetchCmd)
}

// Logic for the fetch CLI subcommand
func fetch(cmd *cobra.Command, args []string) {
	ShowBanner()

	if countRemoteTargets() == 0 {
		core.PrintFatal("no remote values to fetch were specified", 1)
	}

	records := make(map[string]*record.Record, 0)

	// Fetch values from AWs SSM Parameter Store if defined.
	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	if len(awsSsmParameters) != 0 {
		ssmRecords, err := aws.FetchParameterStore()
		if err != nil {
			core.PrintFatal("failed to get SSM parameters", 1)
		}
		core.PrintVerbose(fmt.Sprintf("Fetched %d SSM parameters", len(ssmRecords)))

		for name, record := range ssmRecords {
			records[name] = record
		}
	}

	envFormat, err := record.RecordsAsEnvFile(records)
	if err != nil {
		core.PrintFatal("failed to format records as env file", 1)
	}
	core.PrintNormal("\n")
	core.PrintAlways(envFormat)
	core.PrintNormal(fmt.Sprintf("\nFetched %d values\n", len(records)))
}

// Count the number of user-defined resources to pull values from.
func countRemoteTargets() int {
	remoteTargetCount := 0

	awsSsmParameters := viper.GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	remoteTargetCount += len(awsSsmParameters)

	return remoteTargetCount
}
