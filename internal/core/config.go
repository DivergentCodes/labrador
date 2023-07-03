package core

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var envPrefix = "LAB"
var configFileName = ".labrador.yaml"
var configFileType = "yaml"

// InitConfigDefaults intializes the default configuration settings for the program.
func InitConfigDefaults() {
	initRootDefaults()
	initFetchDefaults()

	initConfigFile()
	initConfigEnv()
}

// Global configuration options (viper lookup strings).
var (
	OptStr_Debug   = "debug"
	OptStr_OutJSON = "out-json"
	OptStr_Quiet   = "quiet"
	OptStr_Verbose = "verbose"
)

func initRootDefaults() {
	viper.SetDefault(OptStr_Debug, false)
	viper.SetDefault(OptStr_OutJSON, false)
	viper.SetDefault(OptStr_Quiet, false)
	viper.SetDefault(OptStr_Verbose, false)
}

// Fetch configuration options
var (
	OptStr_NoConflict = "no-conflict"
	OptStr_OutFile    = "outfile.path"
	OptStr_FileMode   = "outfile.mode"

	OptStr_AWS_Region            = "aws.region"
	OptStr_AWS_SsmParameterStore = "aws.ssm_param"
	OptStr_AWS_SecretManager     = "aws.sm_secret" //#nosec
)

func initFetchDefaults() {
	viper.SetDefault(OptStr_NoConflict, false)
	viper.SetDefault(OptStr_OutFile, "")
	viper.SetDefault(OptStr_FileMode, "0600")

	viper.SetDefault(OptStr_AWS_Region, nil)
	viper.SetDefault(OptStr_AWS_SsmParameterStore, nil)
	viper.SetDefault(OptStr_AWS_SecretManager, nil)
}

// Configuration file and environment.
func initConfigFile() {

	// Use default config file location.
	viper.AddConfigPath(".")
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found. Ignore error and continue.
		} else {
			// Config file was found but another error was produced.
			os.Exit(1)
		}
	}

}

func initConfigEnv() {
	// Support equivalent environment variables.
	viper.SetEnvPrefix(envPrefix)
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
