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

/*
import (
	"encoding/json"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var once sync.Once
var zapLogger *zap.SugaredLogger

func initZapLogger() *zap.SugaredLogger {
	debug := viper.GetBool(OptStr_Debug)

	debugConfig := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": [
			"stdout",
			"/tmp/logs"
		],
		"errorOutputPaths": [
			"stderr"
		],
		"initialFields": {
		},
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	standardConfig := []byte(`{
		"level": "info",
		"encoding": "json",
		"outputPaths": [
			"/tmp/logs"
		],
		"errorOutputPaths": [
			"stderr"
		],
		"initialFields": {
		},
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	var cfg zap.Config
	var err error
	if debug {
		err = json.Unmarshal(debugConfig, &cfg)
	} else {
		err = json.Unmarshal(standardConfig, &cfg)
	}
	if err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger.Sugar()
}

// GetLogger returns a singleton of a configured zap logger.
func GetLogger() *zap.SugaredLogger {

	once.Do(func() {
		zapLogger = initZapLogger()
	})

	return zapLogger
}
*/
