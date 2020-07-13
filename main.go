package main

import (
	"os"

	build "github.com/dodo/dodo-build/pkg/plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/dodo/dodo-build/pkg/command"
)

func main() {
	if os.Getenv(dodo.MagicCookieKey) == dodo.MagicCookieValue {
		build.RegisterPlugin()
		dodo.ServePlugins()
	} else {
		cmd := command.NewBuildCommand()
		if err := cmd.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
