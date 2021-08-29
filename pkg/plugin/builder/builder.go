package builder

import (
	"github.com/dodo-cli/dodo-buildkit/pkg/docker"
	"github.com/dodo-cli/dodo-buildkit/pkg/image"
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/builder"
)

var _ builder.ImageBuilder = &Builder{}

type Builder struct{}

func New() *Builder {
	return &Builder{}
}

func (p *Builder) Type() plugin.Type {
	return builder.Type
}

func (p *Builder) PluginInfo() (*api.PluginInfo, error) {
	return &api.PluginInfo{Name: "build"}, nil
}

func (p *Builder) CreateImage(config *api.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	c, err := docker.GetDockerClient()
	if err != nil {
		return "", err
	}

	img, err := image.NewImage(c, docker.LoadAuthConfig(), config, stream)
	if err != nil {
		return "", err
	}

	imageID, err := img.Get()
	if err != nil {
		return "", err
	}

	return imageID, nil
}
