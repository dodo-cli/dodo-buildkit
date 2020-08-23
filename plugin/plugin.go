package plugin

import (
	"os"

	"github.com/dodo-cli/dodo-build/pkg/command"
	build "github.com/dodo-cli/dodo-build/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
)

func RunMe() int {
	if os.Getenv(dodo.MagicCookieKey) == dodo.MagicCookieValue {
		dodo.ServePlugins(&build.Configuration{})
		return 0
	} else {
		log.SetDefault(log.New(appconfig.GetLoggerOptions()))
		p := &command.Command{}
		if err := p.GetCobraCommand().Execute(); err != nil {
			return 1
		}
		return 0
	}
}

func IncludeMe() {
	dodo.IncludePlugins(&build.Configuration{}, &command.Command{})
}
