package plugin

import (
	"fmt"

	"github.com/dodo-cli/dodo-build/pkg/image"
	"github.com/dodo-cli/dodo-build/pkg/types"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	dodo "github.com/dodo-cli/dodo-core/pkg/types"
	"github.com/dodo-cli/dodo-docker/pkg/client"
	"github.com/oclaussen/go-gimme/configfiles"
)

type Configuration struct{}

func (p *Configuration) Type() plugin.Type {
	return configuration.Type
}

func (p *Configuration) Init() error {
	return nil
}

func (p *Configuration) UpdateConfiguration(backdrop *dodo.Backdrop) (*dodo.Backdrop, error) {
	backdrops := map[string]*types.Backdrop{}
	configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			d := decoder.New(configFile.Path)
			d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
				"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
			})
			return false
		},
	})

	config, err := findBackdrop(backdrops, backdrop.Name)
	if err != nil {
		return &dodo.Backdrop{}, nil
	}

	c, err := client.GetDockerClient()
	if err != nil {
		return nil, err
	}

	img, err := image.NewImage(c, client.LoadAuthConfig(), config.Build)
	if err != nil {
		return nil, err
	}

	imageID, err := img.Get()
	if err != nil {
		return nil, err
	}

	backdrop.ImageId = imageID

	return backdrop, nil
}

func findBackdrop(backdrops map[string]*types.Backdrop, name string) (*types.Backdrop, error) {
	if result, ok := backdrops[name]; ok {
		return result, nil
	}

	for _, b := range backdrops {
		for _, a := range b.Aliases {
			if a == name {
				return b, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", name)
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
