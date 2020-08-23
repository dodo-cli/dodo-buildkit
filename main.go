package main

import (
	"os"

	"github.com/dodo-cli/dodo-build/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
