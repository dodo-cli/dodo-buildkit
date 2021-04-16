package plugin

import (
	build "github.com/dodo-cli/dodo-build/pkg/plugin"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
)

func RunMe() int {
	dodo.ServePlugins(build.New())
	return 0
}

func IncludeMe() {
	dodo.IncludePlugins(build.New())
}
