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

	"github.com/divergentcodes/labrador/internal/core"
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
