package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dodo-cli/dodo-core/pkg/decoder"
	"github.com/dodo-cli/dodo-core/pkg/types"
)

const ErrSSHAgentFormat types.FormatError = "invalid ssh agent format"

func (agent *SshAgent) FromString(spec string) error {
	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 2:
		agent.Id = values[0]
		agent.IdentityFile = values[1]
	default:
		return fmt.Errorf("%s: %w", spec, ErrSSHAgentFormat)
	}

	return nil
}

func NewSSHAgent() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &SshAgent{}
		return &target, DecodeSSHAgent(&target)
	}
}

func DecodeSSHAgent(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	agent := *(target.(**SshAgent))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"id":   decoder.String(&agent.Id),
			"file": decoder.String(&agent.IdentityFile),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)
			if err := agent.FromString(decoded); err != nil {
				d.Error(err)
			}
		},
	})
}
