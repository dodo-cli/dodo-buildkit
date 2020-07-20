package plugin

import (
	"github.com/dodo/dodo-build/pkg/image"
	"github.com/dodo/dodo-build/pkg/types"
	"github.com/dodo/dodo-docker/pkg/client"
	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/plugin"
	dodo "github.com/oclaussen/dodo/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
)

type Configuration struct{}

func RegisterPlugin() {
	plugin.RegisterPluginServer(
		configuration.PluginType,
		&configuration.Plugin{Impl: &Configuration{}},
	)
}

func (p *Configuration) Init() error {
	return nil
}

func (p *Configuration) UpdateConfiguration(backdrop *dodo.Backdrop) (*dodo.Backdrop, error) {
	backdrops := map[string]*types.Backdrop{}
	_, err := configfiles.GimmeConfigFiles(&configfiles.Options{
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

	if err != nil {
		log.L().Error("error finding config files", "error", err)
	}

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

	imageID, err := img.Get()
	if err != nil {
		return nil, err
	}

	backdrop.ImageId = imageID

	return backdrop, nil
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
