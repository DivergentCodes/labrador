package core

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Always print message, even when --quiet is passed.
func PrintAlways(message string) {
	fmt.Print(message)
}

// Always print message, except when --quiet is passed.
func PrintNormal(message string) {
	if !viper.GetBool(OptStr_Quiet) {
		fmt.Print(message)
	}
}

// Only print message when --verbose or --debug is passed.
func PrintVerbose(message string) {
	if viper.GetBool(OptStr_Verbose) || viper.GetBool(OptStr_Debug) {
		fmt.Print(message)
	}
}

// Only print message when --debug is passed.
func PrintDebug(message string) {
	if viper.GetBool(OptStr_Debug) {
		fmt.Print(message)
	}
}

// Print message and immediately exit with exitCode.
func PrintFatal(message string, exitCode int) {
	if exitCode == 0 {
		exitCode = 1
	}
	fmt.Printf("Error: %s\n", message)
	os.Exit(exitCode)
}
