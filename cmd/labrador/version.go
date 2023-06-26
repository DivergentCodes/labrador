package cmd

import (
	"github.com/spf13/cobra"

	"divergent.codes/labrador/internal/core"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `Print the version`,
	Run:   version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	core.PrintAlways(core.Version)
}
