package configuration

import (
	"github.com/dodo/dodo-build/pkg/client"
	"github.com/dodo/dodo-build/pkg/config"
	"github.com/dodo/dodo-build/pkg/image"
	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/plugin/configuration"
	"github.com/oclaussen/dodo/pkg/types"
)

type Configuration struct{}

func NewPlugin() plugin.Plugin {
	return &configuration.Plugin{Impl: &Configuration{}}
}

func (p *Configuration) GetClientOptions(_ string) (*configuration.ClientOptions, error) {
	return &configuration.ClientOptions{}, nil
}

func (p *Configuration) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	conf, err := config.LoadBackdrop(backdrop.Name)
	if err != nil {
		return nil, err
	}

	c, err := client.GetDockerClient()
	if err != nil {
		return nil, err
	}

	img, err := image.NewImage(c, config.LoadAuthConfig(), conf.Build)
	if err != nil {
		return nil, err
	}

	imageId, err := img.Get()
	if err != nil {
		return nil, err
	}

	backdrop.ImageId = imageId
	return backdrop, nil
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
