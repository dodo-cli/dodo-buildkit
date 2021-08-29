package plugin

import (
	"github.com/dodo-cli/dodo-buildkit/pkg/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
)

func RunMe() int {
	if err := plugin.ServePlugins(builder.New()); err != nil {
		log.L().Error("error serving plugin", "error", err)
	}

	return 0
}

func IncludeMe() {
	plugin.IncludePlugins(builder.New())
}
