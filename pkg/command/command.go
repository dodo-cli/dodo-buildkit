package command

import (
	"github.com/dodo/dodo-build/pkg/client"
	"github.com/dodo/dodo-build/pkg/config"
	"github.com/dodo/dodo-build/pkg/image"
	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "build",
		Short:                 "Build all required images for backdrop without running it",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := config.LoadBackdrop(args[0])
			if err != nil {
				return err
			}

			c, err := client.GetDockerClient()
			if err != nil {
				return err
			}

			img, err := image.NewImage(c, config.LoadAuthConfig(), conf.Build)
			if err != nil {
				return err
			}

			_, err = img.Get()
			return err
		},
	}
}
