/*
Labrador fetches variables and secrets from remote services.

Values are recursively pulled from one or more services, and output
to the terminal or a file.

Labrador is focused on reading and pulling values, not on managing
or writing values. It was created with CI/CD pipelines and network
services in mind.

Usage:

	labrador [command]

Available Commands:

	completion  Generate the autocompletion script for the specified shell
	help        Help about any command
	serve       Serve the web API
	version     Print the version

Flags:

	    --config string   config file (default is .app.conf.yaml)
	    --debug           Enable debug mode
	-h, --help            help for labrador
	    --json            Use JSON output
	-q, --quiet           Quiet CLI output
	    --verbose         Verbose CLI output
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"divergent.codes/labrador/internal/core"
)

var (
	rootCmd = &cobra.Command{
		Use:   "labrador",
		Short: "Fetch and load variables and secrets from remote services",
		Long:  "Fetch and load variables and secrets from remote services",
	}
)

// Run for root command. Will execute for all subcommands.
func init() {
	initRootFlags()
}

func initRootFlags() {

	// config
	// TODO: implement reading from configuration file.
	/*
		var cfgFile string
		rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .labrador.conf.yaml)")
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		}
	*/

	// debug
	defaultDebug := viper.GetBool(core.OptStr_Debug)
	rootCmd.PersistentFlags().Bool("debug", defaultDebug, "Enable debug mode")
	err := viper.BindPFlag(core.OptStr_Debug, rootCmd.PersistentFlags().Lookup("debug"))
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
}

// Execute starts the app CLI.
func Execute() error {
	return rootCmd.Execute()
}

// ShowBanner shows the CLI banner message.
func ShowBanner() {

	core.PrintNormal(fmt.Sprintf("Labrador %s created by %s <%s>\n", core.Version, core.AuthorName, core.AuthorEmail))

	core.PrintDebug("Labrador DEBUG mode is enabled\n")

	msg := fmt.Sprintf(
		"\tVersion: %s\n\tCommit: %s\n\tBuilt on: %s\n\tBuilt by: %s\n",
		core.Version, core.Commit, core.Date, core.BuiltBy,
	)
	core.PrintDebug(msg)
}
