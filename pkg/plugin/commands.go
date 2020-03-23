package plugin

import (
	"github.com/dodo/dodo-build/pkg/client"
	"github.com/dodo/dodo-build/pkg/image"
	"github.com/dodo/dodo-build/pkg/types"
	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/spf13/cobra"
)

type Command struct {
	cmds map[string]*cobra.Command
}

func (p *Command) GetCommands() (map[string]*cobra.Command, error) {
	if len(p.cmds) == 0 {
		p.cmds = map[string]*cobra.Command{"build": NewBuildCommand()}
	}
	return p.cmds, nil
}

func NewBuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "build",
		Short:                 "Build all required images for backdrop without running it",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return nil
		},
	}
}
