package builder

import (
	"fmt"

	docker "github.com/docker/docker/client"
	"github.com/dodo-cli/dodo-buildkit/pkg/client"
	"github.com/dodo-cli/dodo-buildkit/pkg/image"
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/builder"
)

var _ builder.ImageBuilder = &Builder{}

type Builder struct {
	client *docker.Client
}

func New() *Builder {
	return &Builder{}
}

func NewFromClient(client *docker.Client) *Builder {
	return &Builder{client: client}
}

func (p *Builder) Type() plugin.Type {
	return builder.Type
}

func (p *Builder) PluginInfo() (*api.PluginInfo, error) {
	return &api.PluginInfo{Name: "buildkit"}, nil
}

func (p *Builder) Client() (*docker.Client, error) {
	if p.client == nil {
		dockerClient, err := client.GetDockerClient()
		if err != nil {
			return nil, fmt.Errorf("could not get docker config: %w", err)
		}

		p.client = dockerClient
	}

	return p.client, nil
}

func (p *Builder) CreateImage(config *api.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	c, err := p.Client()
	if err != nil {
		return "", err
	}

	img, err := image.NewImage(c, client.LoadAuthConfig(), config, stream)
	if err != nil {
		return "", fmt.Errorf("could not initialize builder client: %w", err)
	}

	imageID, err := img.Get()
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return imageID, nil
}
