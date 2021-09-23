package plugin

import (
	"github.com/dodo-cli/dodo-buildkit/pkg/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
)

func RunMe() int {
	m := plugin.Init()
	m.ServePlugins(builder.New())

	return 0
}

func IncludeMe(m plugin.Manager) {
	m.IncludePlugins(builder.New())
}
