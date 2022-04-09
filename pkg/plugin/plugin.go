package plugin

import (
	"github.com/docker/docker/client"
	impl "github.com/dodo-cli/dodo-buildkit/internal/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/builder"
)

func RunMe() int {
	m := plugin.Init()
	m.ServePlugins(NewImageBuilder())

	return 0
}

func IncludeMe(m plugin.Manager) {
	m.IncludePlugins(NewImageBuilder())
}

func NewImageBuilder() builder.ImageBuilder {
	return impl.New()
}

func NewImageBuilderWithDockerClient(c *client.Client) builder.ImageBuilder {
	return impl.NewFromClient(c)
}
