package plugin

import (
	buildkit "github.com/dodo-cli/dodo-buildkit/pkg/plugin"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
)

func RunMe() int {
	dodo.ServePlugins(buildkit.New())
	return 0
}

func IncludeMe() {
	dodo.IncludePlugins(buildkit.New())
}
