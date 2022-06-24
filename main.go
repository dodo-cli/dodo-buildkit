package main

import (
	"os"

	"github.com/wabenet/dodo-buildkit/pkg/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
