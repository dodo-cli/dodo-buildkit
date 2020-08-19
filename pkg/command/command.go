package command

import (
	"github.com/dodo-cli/dodo-build/pkg/image"
	"github.com/dodo-cli/dodo-build/pkg/types"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	"github.com/dodo-cli/dodo-docker/pkg/client"
	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/spf13/cobra"
)

const name = "build"

type Command struct {
	cmd *cobra.Command
}

func (p *Command) Type() plugin.Type {
	return command.Type
}

func (p *Command) Init() error {
	p.cmd = NewBuildCommand()
	return nil
}

func (p *Command) Name() string {
	return name
}

func (p *Command) GetCobraCommand() *cobra.Command {
	return p.cmd
}

func NewBuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   name,
		Short:                 "Build all required images for backdrop without running it",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			for _, backdrop := range backdrops {
				if backdrop.Build != nil && backdrop.Build.ImageName == args[0] {
					c, err := client.GetDockerClient()
					if err != nil {
						return err
					}

					img, err := image.NewImage(c, client.LoadAuthConfig(), backdrop.Build)
					if err != nil {
						return err
					}

					_, err = img.Get()
					return err
				}
			}
			// TODO: error when image config doesn't exist
			return nil
		},
	}
}
