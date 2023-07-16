package core

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// InitConfigDefaults intializes the default configuration settings for the program.
func InitConfigDefaults() {
	initRootDefaults()
	initFetchDefaults()

}

func InitConfigInstance() {
	initConfigFile()
	initConfigEnv()
}

// Global configuration options (viper lookup strings).
var (
	OptStr_Config  = "config"
	OptStr_Debug   = "debug"
	OptStr_OutJSON = "out-json"
	OptStr_Quiet   = "quiet"
	OptStr_Verbose = "verbose"
)

func initRootDefaults() {
	viper.SetDefault(OptStr_Config, "")
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
	OptStr_Quote      = "quote"

	OptStr_AWS_Region            = "aws.region"
	OptStr_AWS_SsmParameterStore = "aws.ssm_param"
	OptStr_AWS_SecretManager     = "aws.sm_secret" //#nosec
)

func initFetchDefaults() {
	viper.SetDefault(OptStr_NoConflict, false)
	viper.SetDefault(OptStr_OutFile, "")
	viper.SetDefault(OptStr_FileMode, "0600")
	viper.SetDefault(OptStr_Quote, false)

	viper.SetDefault(OptStr_AWS_Region, nil)
	viper.SetDefault(OptStr_AWS_SsmParameterStore, nil)
	viper.SetDefault(OptStr_AWS_SecretManager, nil)
}

// Configuration file instance setup.
func initConfigFile() {

	defaultConfigFile := ".labrador.yaml"
	useExplicitConfigFile := true

	// Check if explicit configuration file was passed.
	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = defaultConfigFile
		useExplicitConfigFile = false

	}

	// Config file lookup paths.
	viper.AddConfigPath(".") // Look in current path.
	viper.AddConfigPath("/") // Look in root path so absolute paths work.
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found. Ignore error and continue,
			// unless explicit config file was passed.
			if useExplicitConfigFile {
				PrintFatal(fmt.Sprintf("could not find config file %s", configFile), 1)
			}
		} else {
			// Config file was found but another error was produced.
			PrintFatal(fmt.Sprintf("failed to parse config file %s", configFile), 1)
		}
	}
}

// Environment variable instance setup.
func initConfigEnv() {
	// Support equivalent environment variables.
	viper.SetEnvPrefix("LAB")
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
