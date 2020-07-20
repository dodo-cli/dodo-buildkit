package types

import (
	"fmt"
	"os"
	"strings"

	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

const ErrArgumentFormat types.FormatError = "invalid argument format"

func (arg *Argument) FromString(spec string) error {
	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	case 1:
		arg.Key = values[0]
		arg.Value = os.Getenv(values[0])
	case 2:
		arg.Key = values[0]
		arg.Value = values[1]
	default:
		return fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	}

	return nil
}

func NewArgument() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &Argument{}
		return &target, DecodeArgument(&target)
	}
}

func DecodeArgument(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	arg := *(target.(**Argument))

	return func(d *decoder.Decoder, config interface{}) {
		var decoded string

		decoder.String(&decoded)(d, config)

		if err := arg.FromString(decoded); err != nil {
			d.Error(err)
		}
	}
}
