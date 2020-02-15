package decoder

import (
	"reflect"

	"github.com/dodo/dodo-build/pkg/types"
)

func (d *decoder) DecodeBackdrops(name string, config interface{}) (map[string]types.Backdrop, error) {
	result := map[string]types.Backdrop{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			key := k.(string)
			decoded, err := d.DecodeBackdrop(key, v)
			if err != nil {
				return result, err
			}
			result[key] = decoded
			for _, alias := range decoded.Aliases {
				result[alias] = decoded
			}
		}
	}
	return result, nil
}

func (d *decoder) DecodeBackdrop(name string, config interface{}) (types.Backdrop, error) {
	result := types.Backdrop{Name: name}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "alias", "aliases":
				decoded, err := d.DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.Aliases = decoded
			case "build", "image":
				decoded, err := d.DecodeImage(key, v)
				if err != nil {
					return result, err
				}
				result.Build = &decoded
			}
		}
	}
	return result, nil
}
