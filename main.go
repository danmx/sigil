package main

import (
	"os"

	"github.com/danmx/sigil/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1) //nolint:gomnd
	}
}
