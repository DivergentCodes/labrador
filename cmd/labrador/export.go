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

// Initialize the export CLI subcommand
func init() {
	rootCmd.AddCommand(exportCmd)
}

// Top level logic for the export CLI subcommand
func export(cmd *cobra.Command, args []string) {
	// export implies --quiet
	viper.Set(core.OptStr_Quiet, true)

	variables := fetchVariables()

	toLower := viper.GetBool(core.OptStr_ToLower)
	toUpper := viper.GetBool(core.OptStr_ToUpper)
	formattedOutput, err := variable.VariablesAsShellExport(variables, toLower, toUpper)
	if err != nil {
		core.PrintFatal("failed to format variables as shell exports", 1)
	}

	// Display formatted results to STDOUT.
	core.PrintAlways(formattedOutput)
}
