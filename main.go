package main

import (
	"github.com/bunkr-cli/bunkr/cmd"
	"os"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
