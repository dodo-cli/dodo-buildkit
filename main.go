package main

import (
	"os"

	"github.com/dodo-cli/dodo-buildkit/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
