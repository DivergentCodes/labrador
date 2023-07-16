package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

	rootCmd.AddCommand(fetchCmd)
}

// Top level logic for the fetch CLI subcommand
func fetch(cmd *cobra.Command, args []string) {
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

// Convert the list of variables to formatted output.
func formatVariablesOutput(variables map[string]*variable.Variable) string {
	var formattedOutput string
	var err error

	useQuotes := viper.GetBool(core.OptStr_Quote)
	toLower := viper.GetBool(core.OptStr_ToLower)
	toUpper := viper.GetBool(core.OptStr_ToUpper)

	formattedOutput, err = variable.VariablesAsEnvFile(variables, useQuotes, toLower, toUpper)
	if err != nil {
		core.PrintFatal("failed to format variables as env file", 1)
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
