package main

import (
	"os"

	"github.com/dodo/dodo-build/pkg/command"
	build "github.com/dodo/dodo-build/pkg/plugin"
	dodo "github.com/dodo/dodo-core/pkg/plugin"
)

func main() {
	if os.Getenv(dodo.MagicCookieKey) == dodo.MagicCookieValue {
		dodo.ServePlugins(&build.Configuration{})
	} else {
		cmd := command.NewBuildCommand()
		if err := cmd.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
