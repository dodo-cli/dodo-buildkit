package types

import (
	"os"
	"strings"

	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewArgument() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &Argument{}
		return &target, DecodeArgument(&target)
	}
}

func DecodeArgument(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	env := *(target.(**Argument))
	return func(d *decoder.Decoder, config interface{}) {
		var decoded string
		decoder.String(&decoded)(d, config)
		switch values := strings.SplitN(decoded, "=", 2); len(values) {
		case 1:
			env.Key = values[0]
			env.Value = os.Getenv(values[0])
		case 2:
			env.Key = values[0]
			env.Value = values[1]
		default:
			d.Error("invalid argument")
		}
	}
}
