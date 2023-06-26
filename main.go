package main

import (
	cmd "github.com/divergentcodes/labrador/cmd/labrador"
	"github.com/divergentcodes/labrador/internal/core"
)

func init() {
	// Always start with configured defaults.
	core.InitConfigDefaults()
}

func main() {
	runCli()
}

func runCli() {

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
