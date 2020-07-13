package types

import (
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewSSHAgent() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := SshAgent{}
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
			switch values := strings.SplitN(decoded, "=", 2); len(values) {
			case 2:
				agent.Id = values[0]
				agent.IdentityFile = values[1]
			default:
				d.Error("invalid device")
				return
			}
		},
	})
}
