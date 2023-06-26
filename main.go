package main

import (
	cmd "divergent.codes/labrador/cmd/labrador"
	"divergent.codes/labrador/internal/core"
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
