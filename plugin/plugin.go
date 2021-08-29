package plugin

import (
	"github.com/dodo-cli/dodo-buildkit/pkg/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
)

func RunMe() int {
	plugin.ServePlugins(builder.New())
	return 0
}

func IncludeMe() {
	plugin.IncludePlugins(builder.New())
}
