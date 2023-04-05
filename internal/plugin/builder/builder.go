package builder

import (
	"context"
	"fmt"

	docker "github.com/docker/docker/client"
	"github.com/wabenet/dodo-buildkit/internal/image"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	"github.com/wabenet/dodo-docker/pkg/client"
)

const name = "buildkit"

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

func (p *Builder) PluginInfo() *core.PluginInfo {
	return &core.PluginInfo{
		Name: &core.PluginName{
			Name: name,
			Type: builder.Type.String(),
		},
	}
}

func (p *Builder) Init() (plugin.Config, error) {
	client, err := p.ensureClient()
	if err != nil {
		return nil, err
	}

	ping, err := client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not reach docker host: %w", err)
	}

	return map[string]string{
		"client_version":  client.ClientVersion(),
		"host":            client.DaemonHost(),
		"api_version":     ping.APIVersion,
		"builder_version": fmt.Sprintf("%v", ping.BuilderVersion),
		"os_type":         ping.OSType,
		"experimental":    fmt.Sprintf("%t", ping.Experimental),
	}, nil
}

func (*Builder) Cleanup() {}

func (p *Builder) ensureClient() (*docker.Client, error) {
	if p.client == nil {
		dockerClient, err := client.GetDockerClient()
		if err != nil {
			return nil, fmt.Errorf("could not get docker config: %w", err)
		}

		p.client = dockerClient
	}

	return p.client, nil
}

func (p *Builder) CreateImage(config *core.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	c, err := p.ensureClient()
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
