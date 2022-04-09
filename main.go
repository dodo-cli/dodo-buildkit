package main

import (
	"os"

	"github.com/dodo-cli/dodo-buildkit/pkg/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
