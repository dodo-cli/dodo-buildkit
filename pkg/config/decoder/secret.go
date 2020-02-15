package decoder

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"

	"github.com/dodo/dodo-build/pkg/types"
)

func (d *decoder) DecodeSecrets(name string, config interface{}) ([]*types.Secret, error) {
	result := []*types.Secret{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String, reflect.Map:
		decoded, err := d.DecodeSecret(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeSecret(name, v)
			if err != nil {
				return result, err
			}
			result = append(result, decoded)
		}
	}
	return result, nil
}

func (d *decoder) DecodeSecret(name string, config interface{}) (*types.Secret, error) {
	var result types.Secret
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, t.String())
		if err != nil {
			return nil, err
		}

		reader := csv.NewReader(strings.NewReader(decoded))
		fields, err := reader.Read()
		if err != nil {
			return nil, err
		}

		secretMap := make(map[interface{}]interface{}, len(fields))
		for _, field := range fields {
			dec, err := d.DecodeString(name, field)
			if err != nil {
				return nil, err
			}
			switch values := strings.SplitN(dec, "=", 2); len(values) {
			case 0:
				return nil, fmt.Errorf("empty assignment in '%s'", name)
			case 1:
				return nil, fmt.Errorf("empty assignment in '%s'", name)
			case 2:
				secretMap[values[0]] = values[1]
			default:
				return nil, fmt.Errorf("too many values in '%s'", name)
			}
		}
		return d.DecodeSecret(name, secretMap)
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "id":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Id = decoded
			case "source", "src":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Path = decoded
			}
		}
	}
	return &result, nil
}
