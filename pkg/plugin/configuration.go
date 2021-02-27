package plugin

import (
	"fmt"

	"github.com/dodo-cli/dodo-build/pkg/image"
	"github.com/dodo-cli/dodo-build/pkg/types"
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-docker/pkg/client"
	"github.com/oclaussen/go-gimme/configfiles"
)

var _ configuration.Configuration = &Configuration{}

type Configuration struct{}

func New() *Configuration {
	return &Configuration{}
}

func (p *Configuration) Type() plugin.Type {
	return configuration.Type
}

func (p *Configuration) PluginInfo() (*api.PluginInfo, error) {
	return &api.PluginInfo{Name: "build"}, nil
}

func (p *Configuration) GetBackdrop(alias string) (*api.Backdrop, error) {
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

	config, err := findBackdrop(backdrops, alias)
	if err != nil {
		return &api.Backdrop{}, nil
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

	return &api.Backdrop{ImageId: imageID}, nil
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

func (p *Configuration) ListBackdrops() ([]*api.Backdrop, error) {
	return []*api.Backdrop{}, nil // TODO: implement list
}
