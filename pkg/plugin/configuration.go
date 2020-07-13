package plugin

import (
	"github.com/dodo/dodo-build/pkg/image"
	"github.com/dodo/dodo-build/pkg/types"
	"github.com/dodo/dodo-docker/pkg/client"
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/decoder"
	dodo "github.com/oclaussen/dodo/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
)

type Configuration struct{}

func (p *Configuration) GetClientOptions(_ string) (*configuration.ClientOptions, error) {
	return &configuration.ClientOptions{}, nil
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

	// TODO: wtf this cast
	config, ok := backdrops[backdrop.Name]
	if !ok {
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
