package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/divergentcodes/labrador/internal/core"
	"github.com/divergentcodes/labrador/internal/variable"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Fetch and export values as shell environment variables",
	Long:  "Fetch and export values as shell environment variables",
	Run:   export,
}

// Initialize the fetch CLI subcommand
func init() {

	// aws-region
	defaultAwsRegion := ""
	exportCmd.PersistentFlags().String("aws-region", defaultAwsRegion, "AWS region")
	err := viper.BindPFlag(core.OptStr_AWS_Region, exportCmd.PersistentFlags().Lookup("aws-region"))
	if err != nil {
		panic(err)
	}

	// aws-param
	defaultAwsSsmParameters := viper.GetViper().GetStringSlice(core.OptStr_AWS_SsmParameterStore)
	exportCmd.PersistentFlags().StringSlice("aws-param", defaultAwsSsmParameters, "AWS SSM parameter store path prefix")
	err = viper.BindPFlag(core.OptStr_AWS_SsmParameterStore, exportCmd.PersistentFlags().Lookup("aws-param"))
	if err != nil {
		panic(err)
	}

	// aws-secret
	defaultAwsSmSecrets := viper.GetViper().GetStringSlice(core.OptStr_AWS_SecretManager)
	exportCmd.PersistentFlags().StringSlice("aws-secret", defaultAwsSmSecrets, "AWS Secrets Manager secret name")
	err = viper.BindPFlag(core.OptStr_AWS_SecretManager, exportCmd.PersistentFlags().Lookup("aws-secret"))
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(exportCmd)
}

// Top level logic for the export CLI subcommand
func export(cmd *cobra.Command, args []string) {
	// export implies --quiet
	viper.Set(core.OptStr_Quiet, true)

	if countRemoteTargets() == 0 {
		core.PrintFatal("no remote values to fetch were specified", 1)
	}

	variables := make(map[string]*variable.Variable, 0)
	variables = fetchAwsSsmParameters(variables)
	variables = fetchAwsSmSecrets(variables)

	formattedOutput, err := variable.VariablesAsShellExport(variables)
	if err != nil {
		core.PrintFatal("failed to format variables as shell exports", 1)
	}

	// Display formatted results to STDOUT.
	core.PrintAlways(formattedOutput)
}
