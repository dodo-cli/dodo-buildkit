package command

import (
	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/plugin/command"
	"github.com/spf13/cobra"
)

type Command struct {
	cmd *cobra.Command
}

func NewPlugin() plugin.Plugin {
	return &command.Plugin{Impl: &Command{}}
}

func (p *Command) GetCommand() (*cobra.Command, error) {
	if p.cmd == nil {
		p.cmd = NewBuildCommand()
	}
	return p.cmd, nil
}
